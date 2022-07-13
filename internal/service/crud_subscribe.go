package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"

	"github.com/MlPablo/CRUDService/internal/models"
	"github.com/MlPablo/CRUDService/internal/store"
	"github.com/MlPablo/CRUDService/voc"
)

type SubscribedCRUDService interface {
	SubscribeCreate(context.Context) error
	SubscribeDelete(context.Context) error
	SubscribeUpdate(context.Context) error
	SubscribeGet(context.Context) error
	CRUDService
}

type subscribeCRUDService struct {
	crudService
	cs crudService
	nc *nats.Conn
}

func SubscribeAllCRUD(ctx context.Context, s SubscribedCRUDService) error {
	if err := s.SubscribeCreate(ctx); err != nil {
		return err
	}
	if err := s.SubscribeDelete(ctx); err != nil {
		return err
	}
	if err := s.SubscribeUpdate(ctx); err != nil {
		return err
	}
	if err := s.SubscribeGet(ctx); err != nil {
		return err
	}
	return nil
}

func NewSubscribedService(store store.Storage) (SubscribedCRUDService, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}
	return &subscribeCRUDService{cs: crudService{storage: store}, nc: nc}, nil
}

func (s *subscribeCRUDService) SubscribeCreate(ctx context.Context) error {
	if _, err := s.nc.QueueSubscribe(voc.SubjectCreate, voc.NatsToCrudServicesQueue, func(msg *nats.Msg) {
		user := models.User{}
		if err := json.Unmarshal(msg.Data, &user); err != nil {
			msg.Respond([]byte(err.Error()))
		}

		if err := s.cs.CreateUser(ctx, user); err != nil {
			msg.Respond([]byte(err.Error()))
		}
		log.Println(ctx.Value("name"))
		msg.Respond([]byte(""))
	}); err != nil {
		return err
	}

	return nil
}

func (s *subscribeCRUDService) SubscribeDelete(ctx context.Context) error {
	if _, err := s.nc.QueueSubscribe(voc.SubjectDelete, voc.NatsToCrudServicesQueue, func(msg *nats.Msg) {
		if err := s.cs.DeleteUser(ctx, string(msg.Data)); err != nil {
			msg.Respond([]byte(err.Error()))
		}
		log.Printf("Name of service is %s", ctx.Value("name"))
		msg.Respond([]byte(""))
	}); err != nil {
		return err
	}

	return nil
}

func (s *subscribeCRUDService) SubscribeUpdate(ctx context.Context) error {
	if _, err := s.nc.QueueSubscribe(voc.SubjectUpdate, voc.NatsToCrudServicesQueue, func(msg *nats.Msg) {
		user := models.User{}
		if err := json.Unmarshal(msg.Data, &user); err != nil {
			msg.Respond([]byte(err.Error()))
		}
		if err := s.cs.UpdateUser(ctx, user); err != nil {
			msg.Respond([]byte(err.Error()))
		}
		log.Println(ctx.Value("name"))
		msg.Respond([]byte(""))
	}); err != nil {
		return err
	}

	return nil
}

func (s subscribeCRUDService) SubscribeGet(ctx context.Context) error {
	if _, err := s.nc.QueueSubscribe(voc.SubjectGet, voc.NatsToCrudServicesQueue, func(msg *nats.Msg) {
		user, err := s.cs.GetUser(ctx, string(msg.Data))
		if err == nil {
			log.Println(ctx.Value("name"))
			msg.Respond([]byte(user))
		}
	}); err != nil {
		return err
	}
	return nil
}
