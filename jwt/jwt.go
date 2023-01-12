package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	jwtpkg "github.com/golang-jwt/jwt"
	"github.com/wenccc/myskeleton/app"
	"github.com/wenccc/myskeleton/config"
	"strings"
	"time"
)

var (
	ErrTokenExpired           error = errors.New("令牌已过期")
	ErrTokenExpiredMaxRefresh error = errors.New("令牌已过最大刷新时间")
	ErrTokenMalformed         error = errors.New("请求令牌格式有误")
	ErrTokenInvalid           error = errors.New("请求令牌无效")
	ErrHeaderEmpty            error = errors.New("需要认证才能访问！")
	ErrHeaderMalformed        error = errors.New("请求头中 Authorization 格式有误")
)

type JWT struct {

	// 秘钥，用以加密 JWT，读取配置信息 app.key
	SignKey []byte

	// 刷新 Token 的最大过期时间
	MaxRefresh time.Duration
}

// CustomClaims  负载
type CustomClaims struct {
	UserID       string                 `json:"user_id"`
	UserName     string                 `json:"user_name"`
	ExpireAtTime int64                  `json:"expire_time"`
	OtherData    map[string]interface{} `json:"other_data"` //其他的内容

	// StandardClaims 结构体实现了 Claims 接口继承了  Valid() 方法
	// JWT 规定了7个官方字段，提供使用:
	// - iss (issuer)：发布者
	// - sub (subject)：主题
	// - iat (Issued At)：生成签名的时间
	// - exp (expiration time)：签名过期时间
	// - aud (audience)：观众，相当于接受者
	// - nbf (Not Before)：生效时间
	// - jti (JWT ID)：编号
	jwtpkg.StandardClaims
}

func NewJwt() *JWT {
	return &JWT{
		SignKey:    []byte(config.GetString("app.key")),
		MaxRefresh: time.Duration(config.GetInt64("jwt.max_refresh_time")) * time.Minute,
	}
}

func (jwt *JWT) getTokenFromHeader(ctx *gin.Context) (realToken string, err error) {
	token := ctx.Request.Header.Get("Authorization")
	if len(token) == 0 {
		return "", ErrHeaderEmpty
	}
	p := strings.Split(token, "  ")
	if len(p) != 2 || p[0] != "Bearer" {
		return "", ErrTokenMalformed
	}
	return p[1], nil
}

func (jwt *JWT) parseTokenString(tokenString string) (*jwtpkg.Token, error) {
	return jwtpkg.ParseWithClaims(tokenString, new(CustomClaims), func(token *jwtpkg.Token) (interface{}, error) {
		return jwt.SignKey, nil
	})
}

// expireAtTime 返回过期时间
func (jwt *JWT) expireAtTime() int64 {
	now := app.TimeNowInTimezone()

	expireMin := 0
	if app.IsDebug() {
		expireMin = config.GetInt("jwt.max_refresh_time_debug")
	} else {
		expireMin = config.GetInt("jwt.max_refresh_time")
	}

	return now.Add(time.Duration(expireMin) * time.Minute).Unix()
}

// GenToken 生成token
func (jwt *JWT) GenToken(userID string, userName string, otherData map[string]interface{}) (string, error) {

	claims := CustomClaims{
		UserID:         userID,
		UserName:       userName,
		ExpireAtTime:   jwt.expireAtTime(),
		OtherData:      otherData,
		StandardClaims: jwtpkg.StandardClaims{},
	}

	return jwtpkg.NewWithClaims(jwtpkg.SigningMethodHS256, claims).SignedString(jwt.SignKey)
}

func (jwt *JWT) ParserToken(ctx *gin.Context) (*CustomClaims, error) {

	token, err := jwt.getTokenFromHeader(ctx)
	if err != nil {
		return nil, err
	}

	t, err := jwt.parseTokenString(token)
	if err != nil {
		fmt.Println(err.Error())
	}

	if err != nil {
		validationErr, ok := err.(*jwtpkg.ValidationError)
		if ok {
			if validationErr.Errors == jwtpkg.ValidationErrorMalformed {
				return nil, ErrTokenMalformed
			} else if validationErr.Errors == jwtpkg.ValidationErrorExpired {
				return nil, ErrTokenExpired
			}
		}

		return nil, ErrTokenInvalid
	}
	//这里为啥断言出是指针？
	if c, ok := t.Claims.(*CustomClaims); ok && t.Valid {
		return c, nil
	}

	return nil, ErrTokenInvalid
}

func (jwt *JWT) RefreshToken(c *gin.Context) (string, error) {

	// 1. 从 Header 里获取 token
	tokenString, parseErr := jwt.getTokenFromHeader(c)
	if parseErr != nil {
		return "", parseErr
	}

	// 2. 调用 jwt 库解析用户传参的 Token
	token, err := jwt.parseTokenString(tokenString)

	// 3. 解析出错，未报错证明是合法的 Token（甚至未到过期时间）
	if err != nil {
		validationErr, ok := err.(*jwtpkg.ValidationError)
		// 满足 refresh 的条件：只是单一的报错 ValidationErrorExpired
		if !ok || validationErr.Errors != jwtpkg.ValidationErrorExpired {
			return "", err
		}
	}

	// 4. 解析 JWTCustomClaims 的数据
	claims := token.Claims.(*CustomClaims)

	// 5. 检查是否过了『最大允许刷新的时间』
	x := app.TimeNowInTimezone().Add(-jwt.MaxRefresh).Unix()
	if claims.IssuedAt > x {
		// 修改过期时间
		claims.StandardClaims.ExpiresAt = jwt.expireAtTime()
		return jwt.GenToken(claims.UserID, claims.UserName, claims.OtherData)
	}

	return "", ErrTokenExpiredMaxRefresh
}
