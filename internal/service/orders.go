package service

import (
	"context"

	"github.com/MlPablo/CRUDService/internal/models"
	"github.com/MlPablo/CRUDService/internal/store"
)

type OrderService interface {
	AddOrder(ctx context.Context, order models.Order) error
}

type orderService struct {
	storage store.Storage
}

func NewOrderService(store store.Storage) OrderService {
	return &orderService{storage: store}
}

func (o *orderService) AddOrder(ctx context.Context, order models.Order) error {
	if err := o.storage.NewOrders().Create(ctx, order); err != nil {
		return err
	}
	return nil
}
