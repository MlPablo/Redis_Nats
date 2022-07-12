package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"

	"github.com/MlPablo/CRUDService/internal/models"
)

func (s *server) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := models.User{}
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		nc, err := nats.Connect(nats.DefaultURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Cannot connect to nats")
			return
		}

		byteUser, err := json.Marshal(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Cannot connect to nats")
			return
		}
		resp, err := nc.Request("service.create", byteUser, time.Second)
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

		nc, err := nats.Connect(nats.DefaultURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Cannot connect to nats")
			return
		}

		byteUser, err := json.Marshal(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Cannot connect to nats")
			return
		}
		resp, err := nc.Request("service.update", byteUser, time.Second)
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
		nc, err := nats.Connect(nats.DefaultURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Cannot connect to nats")
			return
		}

		resp, err := nc.Request("service.delete", []byte(id), time.Second)
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

		nc, err := nats.Connect(nats.DefaultURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Cannot connect to nats")
			return
		}

		resp, err := nc.Request("service.get", []byte(id), time.Second)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": string(resp.Data)})
	}
}
