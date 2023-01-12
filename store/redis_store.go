package store

import (
	"github.com/wenccc/myskeleton/redis"
	"time"
)

type Redis struct {
	RedisClient *redis.Redis
	Prefix      string
	expireTime  time.Duration
}

func (r *Redis) Set(key, value string) error {
	r.RedisClient.Set(r.Prefix+key, value, r.expireTime)

	return nil
}

func (r *Redis) Get(key string, clear bool) string {
	key = r.Prefix + key
	value := r.RedisClient.Get(key)
	if clear {
		r.RedisClient.Del(key)
	}
	return value
}

func (r *Redis) Verify(key, answer string, clear bool) bool {
	if len(key) == 0 || len(answer) == 0 {
		return false
	}
	res := r.Get(key, false) == answer
	if res && clear {
		r.RedisClient.Del(key)
	}
	return res
}

func NewRedisStore(expireTime time.Duration, db int, prefix string) (*Redis, error) {
	red, err := redis.GetDefaultRedis(db)
	if err != nil {
		return nil, err
	}
	return &Redis{
		RedisClient: red,
		Prefix:      prefix,
		expireTime:  expireTime,
	}, nil
}
