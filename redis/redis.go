package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type MyRedis struct {
	redis *redis.Client
}

func New(addr, pwd string, db int) (*MyRedis, error) {
	redis := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	})

	_, err := redis.Ping(context.Background()).Result()
	if err != nil {
		err = fmt.Errorf("connect to redis error: %v", err)
		return nil, err
	}

	ins := &MyRedis{
		redis: redis,
	}

	return ins, nil
}
