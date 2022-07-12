package service

import (
	"github.com/MlPablo/CRUDService/internal/models"
	"github.com/MlPablo/CRUDService/internal/store"
)

type CRUDService interface {
	CreateUser(user models.User) error
	UpdateUser(user models.User) error
	GetUser(id string) (string, error)
	DeleteUser(user string) error
}

type crudService struct {
	storage store.Storage
}

func NewCRUDService(store store.Storage) CRUDService {
	return &crudService{storage: store}
}

func (c *crudService) CreateUser(user models.User) error {
	if err := c.storage.Create(user); err != nil {
		return err
	}
	return nil
}

func (c *crudService) UpdateUser(user models.User) error {
	if err := c.storage.Update(user); err != nil {
		return err
	}
	return nil
}
func (c *crudService) GetUser(id string) (string, error) {
	user, err := c.storage.Read(id)
	if err != nil {
		return "", err
	}
	return user, nil
}
func (c *crudService) DeleteUser(user string) error {
	if err := c.storage.Delete(user); err != nil {
		return err
	}
	return nil
}
