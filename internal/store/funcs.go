package store

import (
	"context"
	"errors"
	"time"

	"github.com/MlPablo/CRUDService/internal/models"
)

func (s *storage) Create(ctx context.Context, user models.User) error {
	cont, _ := context.WithTimeout(ctx, time.Second)
	if exist, _ := s.storage.Get(cont, user.User).Result(); exist != "" {
		return errors.New("already exists")
	}

	if err := s.storage.Set(context.Background(), user.User, user.Password, 0).Err(); err != nil {
		return err
	}

	return nil
}

func (s *storage) Read(ctx context.Context, id string) (string, error) {
	cont, _ := context.WithTimeout(ctx, time.Second)
	user := s.storage.Get(cont, id)
	if user.Err() != nil {
		return "", user.Err()
	}
	return user.Val(), nil
}

func (s *storage) Update(ctx context.Context, user models.User) error {
	cont, _ := context.WithTimeout(ctx, time.Second)
	if err := s.storage.Get(cont, user.User).Err(); err != nil {
		return errors.New("no such user")
	}
	if err := s.storage.GetSet(context.Background(), user.User, user.Password).Err(); err != nil {
		return err
	}
	return nil
}

func (s *storage) Delete(ctx context.Context, id string) error {
	cont, _ := context.WithTimeout(ctx, time.Second)
	if err := s.storage.GetDel(cont, id).Err(); err != nil {
		return err
	}
	return nil
}
