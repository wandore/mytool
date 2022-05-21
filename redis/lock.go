package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"time"
)

var LockErr = fmt.Errorf("locking by another node or expired lock")

func (r *MyRedis) TryLock(key, value string, ttl time.Duration) (*Locker, error) {
	ctx := context.Background()

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

func (l *Locker) Keep(ttl time.Duration) error {
	ctx := context.Background()

	release := strconv.FormatInt(int64(ttl/time.Millisecond), 10)

	keepLua := redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pexpire", KEYS[1], ARGV[2]) else return 0 end`)

	ok, err := keepLua.Run(ctx, *l.redis, []string{l.key}, l.uuid+l.value, release).Result()
	if err != nil {
		return err
	}
	if ok == 1 {
		return nil
	}

	return LockErr
}

func (l *Locker) Unlock() error {
	ctx := context.Background()

	unlockLua := redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`)

	ok, err := unlockLua.Run(ctx, *l.redis, []string{l.key}, l.uuid+l.value).Result()
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
