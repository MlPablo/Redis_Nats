package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"

	"github.com/MlPablo/CRUDService/internal/service"
	"github.com/MlPablo/CRUDService/internal/store"
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
	if err := service.SubscribeAll(context.WithValue(context.Background(), "name", "Service 1"), service1); err != nil {
		return nil
	}

	if err := service.SubscribeAll(context.WithValue(context.Background(), "name", "Service 2"), service2); err != nil {
		return nil
	}
	if err := server.router.Run(viper.Get("Host_port").(string)); err != nil {
		return err
	}
	return nil
}
