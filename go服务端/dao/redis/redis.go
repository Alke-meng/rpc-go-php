package redis

import (
	"ccgo/settings"
	"fmt"
	"github.com/go-redis/redis"
)

var Rdb *redis.Client

func Init(cfg *settings.RedisConfig) (err error) {
	Rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PollSize,
	})

	_, err = Rdb.Ping().Result()
	return err

}

func Close() {
	_ = Rdb.Close()
}
