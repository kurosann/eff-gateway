package clients

import (
	"github.com/go-redis/redis/v7"
	"sync"
	"time"
)

type RedisClient struct {
	client *redis.Client

	locker  sync.RWMutex
	isClose bool
}

func NewRedisClient(options redis.Options) *RedisClient {
	return &RedisClient{
		client:  redis.NewClient(&options),
		locker:  sync.RWMutex{},
		isClose: false,
	}
}

func (c *RedisClient) Close() {
	c.locker.Lock()
	defer c.locker.Unlock()

	if !c.isClose {
		c.isClose = true
		_ = c.client.Close()
	}
}
func (c *RedisClient) Set(key string, value interface{}, expiresAt time.Duration) {
	c.client.SetNX(key, value, expiresAt)
}
