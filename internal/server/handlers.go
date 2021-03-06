package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"

	"github.com/MlPablo/CRUDService/internal/models"
	"github.com/MlPablo/CRUDService/voc"
)

func (s *server) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := models.User{}
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		byteUser, err := json.Marshal(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Cannot connect to nats")
			return
		}
		resp, err := s.nc.Request(voc.SubjectCreate, byteUser, time.Second)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(resp.Data) != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": string(resp.Data)})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	}
}

func (s *server) UpdateUser() gin.HandlerFunc {
	user := models.User{}
	return func(c *gin.Context) {
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		byteUser, err := json.Marshal(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Cannot connect to nats")
			return
		}
		resp, err := s.nc.Request(voc.SubjectUpdate, byteUser, time.Second)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(resp.Data) != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": string(resp.Data)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	}
}

func (s *server) DeleteUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("user")

		resp, err := s.nc.Request(voc.SubjectDelete, []byte(id), time.Second)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(resp.Data) != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": string(resp.Data)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	}
}

func (s *server) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("user")

		resp, err := s.nc.Request(voc.SubjectGet, []byte(id), time.Second)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": string(resp.Data)})
	}
}

func (s *server) CreateOrder() gin.HandlerFunc {
	js, _ := s.nc.JetStream()
	msgChan := make(chan *nats.Msg)
	js.ChanSubscribe(voc.SubjectStatusCreateOrder, msgChan)

	return func(c *gin.Context) {
		order := models.Order{}
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		byteUser, err := json.Marshal(order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Cannot connect to nats")
			return
		}

		_, err = js.Publish(voc.SubjectCreateOrder, byteUser)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		select {
		case <-time.After(1 * time.Second):
			c.JSON(http.StatusBadRequest, gin.H{"error": "time exceeded"})
			return
		case msg := <-msgChan:
			if len(msg.Data) != 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": string(msg.Data)})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"status": "created"})
		}
	}
}
