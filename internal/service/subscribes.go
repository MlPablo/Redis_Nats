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

type subscribe struct {
	crudService
	cs crudService
	nc *nats.Conn
}

func SubscribeAll(ctx context.Context, s SubscribedCRUDService) error {
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
	return &subscribe{cs: crudService{storage: store}, nc: nc}, nil
}

func (s *subscribe) SubscribeCreate(ctx context.Context) error {
	if _, err := s.nc.QueueSubscribe(voc.SubjectCreate, voc.NatsToServicesQueue, func(msg *nats.Msg) {
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

func (s *subscribe) SubscribeDelete(ctx context.Context) error {
	if _, err := s.nc.QueueSubscribe(voc.SubjectDelete, voc.NatsToServicesQueue, func(msg *nats.Msg) {
		if err := s.cs.DeleteUser(ctx, string(msg.Data)); err != nil {
			msg.Respond([]byte(err.Error()))
		}
		log.Println(ctx.Value("name"))
		msg.Respond([]byte(""))
	}); err != nil {
		return err
	}

	return nil
}

func (s *subscribe) SubscribeUpdate(ctx context.Context) error {
	if _, err := s.nc.QueueSubscribe(voc.SubjectUpdate, voc.NatsToServicesQueue, func(msg *nats.Msg) {
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

func (s subscribe) SubscribeGet(ctx context.Context) error {
	if _, err := s.nc.QueueSubscribe(voc.SubjectGet, voc.NatsToServicesQueue, func(msg *nats.Msg) {
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
