package redis

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

type Counter struct {
	pool *redis.Pool

	keyPrefix string
}

// New creates a redis based counter instance.
func New(pool *redis.Pool) *Counter {
	return &Counter{
		pool:      pool,
		keyPrefix: "b4fun:counter:",
	}
}

func (c Counter) idKey(id string) string { return c.keyPrefix + id }

func (c *Counter) Get(ctx context.Context, id string) (int64, error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	count, err := redis.Int64(conn.Do("GET", c.idKey(id)))
	switch err {
	case nil:
		return count, nil
	case redis.ErrNil:
		return 0, nil
	default:
		return 0, err
	}
}

func (c *Counter) Incr(ctx context.Context, id string) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("INCR", c.idKey(id))
	return err
}
