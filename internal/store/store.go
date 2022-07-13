package store

import (
	"github.com/go-redis/redis/v9"
	"github.com/spf13/viper"
)

type Storage interface {
	NewCRUD() CRUD
	NewOrders() Orders
}

func New() Storage {
	return &storage{redis.NewClient(
		&redis.Options{
			Addr:     viper.Get("redis_port").(string),
			Password: "",
			DB:       0,
		})}
}

type storage struct {
	store *redis.Client
}

func (s *storage) NewCRUD() CRUD {
	return &crud{s}
}

func (s *storage) NewOrders() Orders {
	return &orders{s}
}
