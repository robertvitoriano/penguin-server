package database

import (
	"github.com/redis/go-redis/v9"
)

type RedisDatabase struct {
	Client   *redis.Client
	Addr     string
	Password string
	DB       int
}

func NewRedisDatabase(addr, password string, db int) *RedisDatabase {
	return &RedisDatabase{
		Addr:     addr,
		Password: password,
		DB:       db,
	}
}

func (r *RedisDatabase) Connect() *redis.Client {
	r.Client = redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.DB,
	})
	return r.Client
}
