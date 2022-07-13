package store

import (
	"context"
	"errors"
	"time"

	"github.com/MlPablo/CRUDService/internal/models"
)

type Orders interface {
	Create(ctx context.Context, user models.Order) error
}

type orders struct {
	*storage
}

func (o *orders) Create(ctx context.Context, order models.Order) error {
	cont, _ := context.WithTimeout(ctx, time.Second)
	if exist, _ := o.store.Get(cont, order.Name).Result(); exist != "" {
		return errors.New("already exists")
	}

	if err := o.store.Set(cont, order.Name, order.Type, 0).Err(); err != nil {
		return err
	}
	return nil
}
