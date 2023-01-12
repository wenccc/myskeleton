package app

import (
	"github.com/wenccc/myskeleton/config"
	"time"
)

func IsLocal() bool {

	return config.Get("app.env") == "local"
}

func IsProduction() bool {
	return config.Get("app.env") == "production"
}
func IsDebug() bool {
	return config.GetBool("app.Debug")
}

func IsTesting() bool {
	return config.Get("app.env") == "testing"
}

func TimeNowInTimezone() time.Time {
	chinaTimezone, _ := time.LoadLocation(config.GetString("app.timezone", "Asia/Shanghai"))
	return time.Now().In(chinaTimezone)
}
