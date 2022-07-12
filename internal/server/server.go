package server

import (
	"github.com/gin-gonic/gin"

	"github.com/MlPablo/CRUDService/internal/service"
	"github.com/MlPablo/CRUDService/internal/store"
)

//type Server interface {
//	CreateUser()
//	UpdateUser()
//	DeleteUserByID()
//	GetUserByID()
//}

type server struct {
	router  *gin.Engine
	storage store.Storage
}

func NewServer() *server {
	return &server{router: gin.Default(), storage: store.New()}
}

func Start() error {
	server := NewServer()
	server.SetUpRoutes()

	service1 := service.NewSubscribedService(server.storage)
	service.SubcribeAll(service1, "Service 1")

	service2 := service.NewSubscribedService(server.storage)
	service.SubcribeAll(service2, "Service 2")

	if err := server.router.Run(":8080"); err != nil {
		return err
	}
	return nil
}
