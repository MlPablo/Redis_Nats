package service

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"

	"github.com/MlPablo/CRUDService/internal/models"
	"github.com/MlPablo/CRUDService/internal/store"
)

type SubscribedCRUDService interface {
	SubscribeCreate(string) error
	SubscribeDelete(string) error
	SubscribeUpdate(string) error
	SubscribeGet(string) error
	CRUDService
}

type subscribe struct {
	crudService
	cs crudService
}

func SubcribeAll(s SubscribedCRUDService, name string) error {
	if err := s.SubscribeCreate(name); err != nil {
		return err
	}
	if err := s.SubscribeDelete(name); err != nil {
		return err
	}
	if err := s.SubscribeUpdate(name); err != nil {
		return err
	}
	if err := s.SubscribeGet(name); err != nil {
		return err
	}
	return nil
}

func NewSubscribedService(store store.Storage) SubscribedCRUDService {
	return &subscribe{cs: crudService{storage: store}}
}

func (s *subscribe) SubscribeCreate(name string) error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	if _, err := nc.QueueSubscribe("service.create", "service_queue", func(msg *nats.Msg) {
		user := models.User{}
		if err := json.Unmarshal(msg.Data, &user); err != nil {
			msg.Respond([]byte(err.Error()))
		}
		log.Println(user)
		if err := s.cs.CreateUser(user); err != nil {
			msg.Respond([]byte(err.Error()))
		}
		log.Println(name)
		msg.Respond([]byte(""))
	}); err != nil {
		return err
	}

	return nil
}

func (s *subscribe) SubscribeDelete(name string) error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	if _, err := nc.QueueSubscribe("service.delete", "service_queue", func(msg *nats.Msg) {
		if err := s.cs.DeleteUser(string(msg.Data)); err != nil {
			msg.Respond([]byte(err.Error()))
		}
		log.Println(name)
		msg.Respond([]byte(""))
	}); err != nil {
		return err
	}

	return nil
}

func (s *subscribe) SubscribeUpdate(name string) error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	if _, err := nc.QueueSubscribe("service.update", "service_queue", func(msg *nats.Msg) {
		user := models.User{}
		if err := json.Unmarshal(msg.Data, &user); err != nil {
			msg.Respond([]byte(err.Error()))
		}
		if err := s.cs.UpdateUser(user); err != nil {
			msg.Respond([]byte(err.Error()))
		}
		log.Println(name)
		msg.Respond([]byte(""))
	}); err != nil {
		return err
	}

	return nil
}

func (s subscribe) SubscribeGet(name string) error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	if _, err := nc.QueueSubscribe("service.get", "service_queue", func(msg *nats.Msg) {
		user, err := s.cs.GetUser(string(msg.Data))
		log.Println(name)
		if err == nil {
			msg.Respond([]byte(user))
		}
		//msg.Respond([]byte(user))
	}); err != nil {
		return err
	}

	return nil
}