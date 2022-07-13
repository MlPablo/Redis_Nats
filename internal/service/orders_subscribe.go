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

type SubscribedOrderService interface {
	SubscribeCreate(context.Context) error
	OrderService
}

type subscribeOrderService struct {
	orderService
	os orderService
	nc *nats.Conn
}

func SubscribeAllOrders(ctx context.Context, s SubscribedOrderService) error {
	if err := s.SubscribeCreate(ctx); err != nil {
		return err
	}
	return nil
}

func NewSubscribedOrderService(store store.Storage) (SubscribedOrderService, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}
	return &subscribeOrderService{os: orderService{storage: store}, nc: nc}, nil
}

func (s *subscribeOrderService) SubscribeCreate(ctx context.Context) error {
	js, err := s.nc.JetStream()
	if err != nil {
		return err
	}

	if _, err := js.QueueSubscribe(voc.SubjectCreateOrder, voc.NatsToOrderServicesQueue, func(msg *nats.Msg) {
		msg.Ack()
		order := models.Order{}
		log.Printf("Name of service is %s", ctx.Value("name"))
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			js.Publish(voc.SubjectStatusCreateOrder, []byte(err.Error()))
			return
		}
		if err := s.os.AddOrder(ctx, order); err != nil {
			js.Publish(voc.SubjectStatusCreateOrder, []byte(err.Error()))
			return
		}
		js.Publish(voc.SubjectStatusCreateOrder, []byte(""))
	}); err != nil {
		return err
	}

	return nil
}
