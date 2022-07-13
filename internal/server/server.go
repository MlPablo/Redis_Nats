package server

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"

	"github.com/MlPablo/CRUDService/internal/service"
	"github.com/MlPablo/CRUDService/internal/store"
	"github.com/MlPablo/CRUDService/voc"
)

type server struct {
	router  *gin.Engine
	storage store.Storage
	nc      *nats.Conn
}

func NewServer() (*server, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	return &server{router: gin.Default(), storage: store.New(), nc: nc}, nil
}

func Start() error {
	server, err := NewServer()
	if err != nil {
		return err
	}
	server.SetUpRoutes()

	service1, err := service.NewSubscribedService(server.storage)
	if err != nil {
		return err
	}
	service2, err := service.NewSubscribedService(server.storage)
	if err != nil {
		return err
	}
	if err := service.SubscribeAllCRUD(context.WithValue(context.Background(), "name", "Service 1"), service1); err != nil {
		return nil
	}

	if err := service.SubscribeAllCRUD(context.WithValue(context.Background(), "name", "Service 2"), service2); err != nil {
		return nil
	}

	service3, err := service.NewSubscribedOrderService(server.storage)
	if err != nil {
		return err
	}

	service4, err := service.NewSubscribedOrderService(server.storage)
	if err != nil {
		return err
	}

	js, err := server.nc.JetStream()
	if err != nil {
		return err
	}

	stream, _ := js.StreamInfo(voc.NatsOrderStreamName)
	if stream == nil {
		log.Printf("creating stream %q", voc.NatsOrderStreamName)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     voc.NatsOrderStreamName,
			Subjects: []string{voc.SubjectCreateOrder, voc.SubjectStatusCreateOrder},
			MaxAge:   0,
			Storage:  nats.MemoryStorage,
		})
		if err != nil {
			return err
		}
	}

	if err := service.SubscribeAllOrders(context.WithValue(context.Background(), "name", "Service 3"), service3); err != nil {
		return err
	}

	if err := service.SubscribeAllOrders(context.WithValue(context.Background(), "name", "Service 4"), service4); err != nil {
		return err
	}
	if err := server.router.Run(viper.Get("Host_port").(string)); err != nil {
		return err
	}
	return nil
}
