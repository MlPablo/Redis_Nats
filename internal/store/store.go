package store

import (
	"github.com/go-redis/redis/v9"

	"github.com/MlPablo/CRUDService/internal/models"
)

type Storage interface {
	Create(user models.User) error
	Read(string) (string, error)
	Update(user models.User) error
	Delete(string) error
}

func New() Storage {
	return &storage{redis.NewClient(
		&redis.Options{
			Addr:     ":6379",
			Password: "",
			DB:       0,
		})}
}

type storage struct {
	storage *redis.Client
}
