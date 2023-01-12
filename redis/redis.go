package redis

import (
	"context"
	"fmt"
	goRedis "github.com/go-redis/redis/v8"
	"github.com/wenccc/myskeleton/configcenter"
	"github.com/wenccc/myskeleton/logger"
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	pool = struct {
		lock sync.RWMutex
		pool map[int]*Redis
	}{lock: sync.RWMutex{}, pool: make(map[int]*Redis)}
)

const (
	loggerModule = "redis"
)

type Redis struct {
	Client  *goRedis.Client
	Db      int
	Context context.Context
}

func newRedis(conf configcenter.RedisConf, db ...int) *Redis {
	dbNumber := 0
	if len(db) > 0 {
		dbNumber = db[0]
	}

	opts := &goRedis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Username: conf.UserName,
		Password: conf.Password,
		DB:       dbNumber,
	}
	cc := goRedis.NewClient(opts)
	ping := cc.Ping(context.Background())
	if ping.Err() != nil {
		logger.Error(loggerModule, zap.String("err", ping.Err().Error()))
	}

	return &Redis{
		Client:  cc,
		Db:      dbNumber,
		Context: context.Background(),
	}
}

func GetDefaultRedis(db ...int) (*Redis, error) {
	dbNumber := 0
	if len(db) > 0 {
		dbNumber = db[0]
	}

	redisConf, err := configcenter.GetDefaultRedis()
	if err != nil {
		return nil, err
	}

	pool.lock.RLock()
	cc, ok := pool.pool[dbNumber]
	if !ok {
		pool.lock.RUnlock()
		pool.lock.Lock()
		cc, ok = pool.pool[dbNumber]
		if !ok {
			cc = newRedis(redisConf, dbNumber)
			pool.pool[dbNumber] = cc
		}
		pool.lock.Unlock()
		return cc, nil
	}
	pool.lock.RUnlock()
	return cc, nil
}

func (r Redis) Ping() error {
	return r.Client.Ping(r.Context).Err()
}

func (r Redis) Get(key string) string {
	res, _ := r.Client.Get(r.Context, key).Result()
	return res
}
func (r Redis) Del(key string) bool {

	err := r.Client.Del(r.Context, key).Err()
	if err != nil {
		logger.Error(loggerModule, zap.String("Del", err.Error()), zap.String("key", key))
		return false
	}
	return true
}

func (r Redis) Set(key string, value string, expire ...time.Duration) bool {

	exp := time.Duration(0)
	if len(expire) > 0 {
		exp = expire[0]
	}

	if err := r.Client.Set(r.Context, key, value, exp).Err(); err != nil {
		logger.Error(loggerModule, zap.String("Set", err.Error()), zap.String("key", key), zap.String("value", value))
		return false
	}

	return true
}

func (r Redis) Has(key string) bool {
	res := r.Client.Exists(r.Context, key)

	if res.Err() != nil {
		logger.Error(loggerModule, zap.String("Exists", res.Err().Error()), zap.String("key", key))
		return false
	}
	u, err := res.Uint64()
	if err != nil {
		logger.Error(loggerModule, zap.String("Exists Uint64", res.Err().Error()), zap.String("key", key))
		return false
	}

	return u > 0
}

func (r Redis) Increment(key string, val ...float64) bool {
	var incrementValue float64 = 1
	if len(val) > 0 {
		incrementValue = val[0]
	}
	if err := r.Client.IncrByFloat(r.Context, key, incrementValue).Err(); err != nil {
		logger.Error(loggerModule, zap.String("Increment Uint64", err.Error()), zap.String("key", key), zap.Float64("value", incrementValue))
		return false
	}

	return true
}
