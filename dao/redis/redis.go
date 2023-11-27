package redis

import (
	"SHforum_backend/settings"
	"fmt"
	"github.com/go-redis/redis"
)

var (
	client *redis.Client
	Nil    = redis.Nil
)

func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		Password: cfg.Password, //取不到默认为空
		DB:       cfg.DB,       //默认为0
		PoolSize: cfg.PoolSize,
	})

	_, err = client.Ping().Result()
	return err
}

func Close() {
	_ = client.Close()
}
