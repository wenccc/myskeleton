package configcenter

type MysqlConf struct {
	Host               string
	User               string
	Password           string
	Port               int
	DataBase           string
	MaxIdleConnections int
	MaxLifeSeconds     int
	MaxOpenConnections int
}

var (
	defaultMysqlConf *MysqlConf
)

func SetMysql(conf MysqlConf) {
	defaultMysqlConf = &conf
}

// GetDefaultMysql GetDefaultRedis 获取默认配置
func GetDefaultMysql() (MysqlConf, error) {
	if defaultRedisConf == nil {
		return MysqlConf{}, ErrNotConfig
	}
	return *defaultMysqlConf, nil
}
