package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"

	"github.com/MlPablo/CRUDService/internal/service"
	"github.com/MlPablo/CRUDService/internal/store"
)

const (
	SubjectCreate = "service.create"
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
	if err := service.SubscribeAll(service1, "Service 1"); err != nil {
		return nil
	}

	if err := service.SubscribeAll(service2, "Service 2"); err != nil {
		return nil
	}
	if err := server.router.Run(viper.Get("Host_port").(string)); err != nil {
		return err
	}
	return nil
}
