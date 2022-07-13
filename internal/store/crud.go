package store

import (
	"context"
	"errors"
	"time"

	"github.com/MlPablo/CRUDService/internal/models"
)

type CRUD interface {
	Create(ctx context.Context, user models.User) error
	Read(ctx context.Context, id string) (string, error)
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id string) error
}

type crud struct {
	*storage
}

func (c *crud) Create(ctx context.Context, user models.User) error {
	cont, _ := context.WithTimeout(ctx, time.Second)
	if exist, _ := c.store.Get(cont, user.User).Result(); exist != "" {
		return errors.New("already exists")
	}

	if err := c.store.Set(context.Background(), user.User, user.Password, 0).Err(); err != nil {
		return err
	}

	return nil
}

func (c *crud) Read(ctx context.Context, id string) (string, error) {
	cont, _ := context.WithTimeout(ctx, time.Second)
	user := c.store.Get(cont, id)
	if user.Err() != nil {
		return "", user.Err()
	}
	return user.Val(), nil
}

func (c *crud) Update(ctx context.Context, user models.User) error {
	cont, _ := context.WithTimeout(ctx, time.Second)
	if err := c.store.Get(cont, user.User).Err(); err != nil {
		return errors.New("no such user")
	}
	if err := c.store.GetSet(context.Background(), user.User, user.Password).Err(); err != nil {
		return err
	}
	return nil
}

func (c *crud) Delete(ctx context.Context, id string) error {
	cont, _ := context.WithTimeout(ctx, time.Second)
	if err := c.store.GetDel(cont, id).Err(); err != nil {
		return err
	}
	return nil
}
