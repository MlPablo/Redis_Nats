package store

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/spf13/viper"

	"github.com/MlPablo/CRUDService/internal/models"
)

type Storage interface {
	Create(ctx context.Context, user models.User) error
	Read(ctx context.Context, id string) (string, error)
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id string) error
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
	storage *redis.Client
}
