package service

import (
	"context"

	"github.com/MlPablo/CRUDService/internal/models"
	"github.com/MlPablo/CRUDService/internal/store"
)

type CRUDService interface {
	CreateUser(ctx context.Context, user models.User) error
	UpdateUser(ctx context.Context, user models.User) error
	GetUser(ctx context.Context, id string) (string, error)
	DeleteUser(ctx context.Context, user string) error
}

type crudService struct {
	storage store.Storage
}

func NewCRUDService(store store.Storage) CRUDService {
	return &crudService{storage: store}
}

func (c *crudService) CreateUser(ctx context.Context, user models.User) error {
	if err := c.storage.Create(ctx, user); err != nil {
		return err
	}
	return nil
}

func (c *crudService) UpdateUser(ctx context.Context, user models.User) error {
	if err := c.storage.Update(ctx, user); err != nil {
		return err
	}
	return nil
}
func (c *crudService) GetUser(ctx context.Context, id string) (string, error) {
	user, err := c.storage.Read(ctx, id)
	if err != nil {
		return "", err
	}
	return user, nil
}
func (c *crudService) DeleteUser(ctx context.Context, user string) error {
	if err := c.storage.Delete(ctx, user); err != nil {
		return err
	}
	return nil
}
