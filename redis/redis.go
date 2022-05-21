package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
)

var LockErr = fmt.Errorf("locking by another node")

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

func (r *MyRedis) TryLock(ctx context.Context, key, value string, ttl time.Duration) (*Locker, error) {
	uuid := uuid.NewV4().String()

	ok, err := r.redis.SetNX(ctx, key, uuid+value, ttl).Result()
	if err != nil {
		err = fmt.Errorf("try lock error: %v", err)
		return nil, err
	}
	if !ok {
		return nil, LockErr
	}

	locker := &Locker{
		redis: r.redis,
		key:   key,
		value: value,
		uuid:  uuid,
	}
	return locker, nil
}

type Locker struct {
	redis *redis.Client
	key   string
	value string
	uuid  string
}

func (l *Locker) Keep(ctx context.Context, ttl time.Duration) error {
	keepLua := redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pexpire", KEYS[1], ARGV[2]) else return 0 end`)

	ok, err := keepLua.Run(ctx, l.redis, []string{l.key, l.uuid + l.value}, ttl).Result()
	if err != nil {
		return err
	}
	if ok == 1 {
		return nil
	}

	return LockErr
}

func (l *Locker) Unlock(ctx context.Context) error {
	unlockLua := redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`)

	ok, err := unlockLua.Run(ctx, l.redis, []string{l.key, l.uuid + l.value}).Result()
	if err == redis.Nil {
		return LockErr
	}
	if err != nil {
		return err
	}

	if ok != 1 {
		return LockErr
	}

	return nil
}
