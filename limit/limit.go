package limit

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	limitStoreDriver "github.com/ulule/limiter/v3/drivers/store/redis"
	"github.com/wenccc/myskeleton/config"
	"github.com/wenccc/myskeleton/logger"
	"github.com/wenccc/myskeleton/redis"
	"go.uber.org/zap"
	"strings"
)

func GetKeyIP(ctx *gin.Context) string {

	return ctx.ClientIP()
}

func GetKeyRouteIp(ctx *gin.Context) string {

	return strings.ReplaceAll(ctx.FullPath(), "/", "-") + GetKeyIP(ctx)
}

func CheckRate(ctx *gin.Context, key, formatted string) (limiter.Context, error) {
	var c limiter.Context
	rate, err := limiter.NewRateFromFormatted(formatted)
	if err != nil {
		logger.Fatal("limit", zap.Error(err))
		return c, err
	}

	r, err := redis.GetDefaultRedis(0)
	if err != nil {
		return c, err
	}

	store, err := limitStoreDriver.NewStoreWithOptions(r.Client, limiter.StoreOptions{
		Prefix: config.GetString("app.app_name") + ":limiter:",
	})

	if err != nil {
		logger.Fatal("limit.limitStoreDriver", zap.Error(err))
		return c, err
	}

	l := limiter.New(store, rate)
	if ctx.GetBool("has-limit") {
		return l.Peek(ctx, key)
	}
	ctx.Set("has-limit", true)
	return l.Get(ctx, key)
}
