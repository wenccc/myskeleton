package config

import (
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/wenccc/myskeleton/util"
	"os"
)

var c *viper.Viper

type Func func() map[string]interface{}

var Funcs map[string]Func

func init() {
	c = viper.New()
	c.SetConfigType("env")
	c.AddConfigPath(".")
	c.SetEnvPrefix("APPENV")
	c.AutomaticEnv()

	Funcs = make(map[string]Func)
}

func InitConfig(envSuffix string) {
	LoadEnv(envSuffix)
	loadConfig(Funcs)
}

func loadConfig(funs map[string]Func) {
	for name, fun := range funs {
		c.Set(name, fun())
	}
}

// Add 新增配置项
func Add(name string, configFn Func) {
	Funcs[name] = configFn
}

func LoadEnv(envSuffix string) {
	defaultFile := ".env"
	if len(envSuffix) > 0 {
		defaultFile = defaultFile + "." + envSuffix
		if _, err := os.Stat(defaultFile); err != nil {
			panic(err)
		}
	}
	c.SetConfigName(defaultFile)

	if err := c.ReadInConfig(); err != nil {
		panic(err)
	}
	c.WatchConfig()
}

func internalGet(path string, defaultValue ...interface{}) interface{} {

	if !c.IsSet(path) || util.Empty(c.Get(path)) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}

	return c.Get(path)
}

func Env(path string, defaultValue ...interface{}) interface{} {
	if len(defaultValue) > 0 {
		return internalGet(path, defaultValue[0])
	}

	return internalGet(path)
}

func GetString(path string, defaultValue ...interface{}) string {
	return cast.ToString(internalGet(path, defaultValue...))
}

func Get(path string, defaultValue ...interface{}) string {
	return GetString(path, defaultValue...)
}

func GetInt(path string, defaultValue ...interface{}) int {
	return cast.ToInt(internalGet(path, defaultValue...))
}
func GetInt64(path string, defaultValue ...interface{}) int64 {
	return cast.ToInt64(internalGet(path, defaultValue...))
}
func GetFloat64(path string, defaultValue ...interface{}) float64 {
	return cast.ToFloat64(internalGet(path, defaultValue...))
}
func GetBool(path string, defaultValue ...interface{}) bool {
	return cast.ToBool(internalGet(path, defaultValue...))
}

// GetStringMapString 获取结构数据
func GetStringMapString(path string) map[string]string {
	return c.GetStringMapString(path)
}
