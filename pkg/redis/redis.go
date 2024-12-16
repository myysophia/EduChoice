package redis

import (
	"context"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() (*redis.Client, error) {
	rd := redis.NewClient(&redis.Options{
		Addr:     common.CONFIG.String("redis.host") + ":" + common.CONFIG.String("redis.port"),
		Password: common.CONFIG.String("redis.password"),
		DB:       0,  // 默认DB 1 连接到服务器后要选择的数据库。
		PoolSize: 20, // 最大套接字连接数。 默认情况下，每个可用CPU有10个连接，由runtime.GOMAXPROCS报告。
	})
	ctx := context.Background()
	if _, err := rd.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	return rd, nil
}
