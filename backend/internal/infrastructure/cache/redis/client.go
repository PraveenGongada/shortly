package redis

import "github.com/redis/go-redis/v9"

type Client interface {
	Client() *redis.Client
}

type client struct {
	rdb *redis.Client
}

func (c *client) Client() *redis.Client {
	return c.rdb
}

func NewClient() Client {
	return &client{}
}
