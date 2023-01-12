package configcenter

type RedisConf struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

var defaultRedisConf *RedisConf

func SetRedis(conf RedisConf) {
	defaultRedisConf = &conf
}

// GetDefaultRedis 获取默认redis配置
func GetDefaultRedis() (RedisConf, error) {
	if defaultRedisConf == nil {
		return RedisConf{}, ErrNotConfig
	}
	return *defaultRedisConf, nil
}
