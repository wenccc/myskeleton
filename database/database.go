package database

import (
	"fmt"
	"github.com/wenccc/myskeleton/configcenter"
	log "github.com/wenccc/myskeleton/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var defaultDbConn *gorm.DB

func GetDefaultConn() (*gorm.DB, error) {
	if defaultDbConn != nil {
		return defaultDbConn, nil
	}
	conf, err := configcenter.GetDefaultMysql()
	if err != nil {
		return nil, err
	}
	defaultDbConn, err = NewMysqlConn(conf)
	return defaultDbConn, err
}

func NewMysqlConn(conf configcenter.MysqlConf) (conn *gorm.DB, err error) {

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&multiStatements=true&loc=Local",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DataBase,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: log.NewGormLogger(),
	})
	if err != nil {
		return nil, err
	}

	rawDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	if conf.MaxIdleConnections > 0 {
		rawDb.SetMaxIdleConns(conf.MaxIdleConnections)
	}

	if conf.MaxLifeSeconds > 0 {
		rawDb.SetConnMaxLifetime(time.Duration(conf.MaxLifeSeconds) * time.Second)
	}
	if conf.MaxOpenConnections > 0 {
		rawDb.SetMaxOpenConns(conf.MaxOpenConnections)
	}

	return db, nil
}
