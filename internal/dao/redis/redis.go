package redis

import (
	"SHforum_backend/internal/settings"
	"context"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"sync"
)

var (
	client  *goredislib.Client
	once    sync.Once
	RedSync *redsync.Redsync
	Nil     = goredislib.Nil
	ctx     = context.Background()
)

func Init(cfg *settings.RedisConfig) (err error) {
	once.Do(func() {
		client = goredislib.NewClient(&goredislib.Options{
			Addr: fmt.Sprintf("%s:%d",
				cfg.Host,
				cfg.Port,
			),
			Password: cfg.Password, //取不到默认为空
			DB:       cfg.DB,       //默认为0
			PoolSize: cfg.PoolSize,
		})
	})
	pool := goredis.NewPool(client)
	RedSync = redsync.New(pool)
	_, err = client.Ping(context.TODO()).Result()
	return err
}

func Close() {
	_ = client.Close()
}
