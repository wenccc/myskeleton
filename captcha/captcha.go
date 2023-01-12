package captcha

import (
	"fmt"
	"github.com/mojocn/base64Captcha"
	"github.com/wenccc/myskeleton/app"
	"github.com/wenccc/myskeleton/config"
	"github.com/wenccc/myskeleton/store"
	"sync"
	"time"
)

type Captcha struct {
	Base64Captcha *base64Captcha.Captcha
}

var (
	pool = struct {
		lock      sync.RWMutex
		container map[string]*Captcha
	}{
		lock:      sync.RWMutex{},
		container: make(map[string]*Captcha),
	}

	defaultOption = Option{
		Height:   20,
		Width:    50,
		Length:   20,
		DotCount: 10,
		MaxSkew:  10,
	}
)

type Option struct {
	Height, Width, Length, DotCount int
	MaxSkew                         float64
}

func NewCaptcha(opts ...Option) (*Captcha, error) {

	opt := defaultOption
	if len(opts) > 0 {
		opt = opts[0]
	}
	key := fmt.Sprintf("%d_%d_%d_%d_%f", opt.Height, opt.Width, opt.Length, opt.DotCount, opt.MaxSkew)
	pool.lock.RLock()
	c, ok := pool.container[key]
	if !ok {
		pool.lock.RUnlock()
		pool.lock.Lock()
		c, ok = pool.container[key]
		if !ok {
			digit := base64Captcha.NewDriverDigit(opt.Height, opt.Width, opt.Length, opt.MaxSkew, opt.DotCount)
			c = new(Captcha)

			expireTime := time.Duration(config.GetInt("captcha.expire_time")) * time.Minute

			s, err := store.NewRedisStore(expireTime, config.GetInt("captcha.db", 0), config.GetString("captcha.prefix"))
			if err != nil {
				return nil, err
			}
			c.Base64Captcha = base64Captcha.NewCaptcha(digit, s)
			pool.container[key] = c
			pool.lock.Unlock()

		}
		return c, nil
	}
	pool.lock.RUnlock()

	return c, nil
}

// GenerateCaptcha 生成图片验证码
func (c *Captcha) GenerateCaptcha() (id string, b64s string, err error) {
	return c.Base64Captcha.Generate()
}

// VerifyCaptcha 验证验证码是否正确
func (c *Captcha) VerifyCaptcha(id string, answer string) (match bool) {

	// 方便本地和 API 自动测试
	if !app.IsProduction() && id == config.GetString("captcha.testing_key") {
		return true
	}
	// 第三个参数是验证后是否删除，我们选择 false
	// 这样方便用户多次提交，防止表单提交错误需要多次输入图片验证码
	return c.Base64Captcha.Verify(id, answer, false)
}
