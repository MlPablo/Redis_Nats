package store

import (
	"context"
	"errors"

	"github.com/MlPablo/CRUDService/internal/models"
)

func (s *storage) Create(user models.User) error {
	if exist, _ := s.storage.Get(context.Background(), user.User).Result(); exist != "" {
		return errors.New("already exists")
	}

	if err := s.storage.Set(context.Background(), user.User, user.Password, 0).Err(); err != nil {
		return err
	}

	return nil
}

func (s *storage) Read(id string) (string, error) {
	user := s.storage.Get(context.Background(), id)
	if user.Err() != nil {
		return "", user.Err()
	}
	return user.Val(), nil
}

func (s *storage) Update(user models.User) error {
	if err := s.storage.Get(context.Background(), user.User).Err(); err != nil {
		return errors.New("no such user")
	}
	if err := s.storage.GetSet(context.Background(), user.User, user.Password).Err(); err != nil {
		return err
	}
	return nil
}

func (s *storage) Delete(id string) error {
	if err := s.storage.GetDel(context.Background(), id).Err(); err != nil {
		return err
	}
	return nil
}
