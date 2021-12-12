package freestuff

import "github.com/gomodule/redigo/redis"

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

var pool = newPool()

type LinkCache interface {
	IsKnown(title string) bool
	SetKnown(title string) error
}

type RedisCache struct {
	conn redis.Conn
}

func NewRedisCache() RedisCache {
	conn := pool.Get()
	return RedisCache{conn: conn}
}

func (r *RedisCache) IsKnown(title string) (bool, error) {
	value, err := r.conn.Do("EXISTS", title)
	if err != nil {
		return false, err
	}
	return value == int64(1), nil
}

func (r *RedisCache) SetKnown(title string) error {
	const TTL = "1209600"
	_, err := r.conn.Do("SETEX", title, TTL, "")
	return err
}
