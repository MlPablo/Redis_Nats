package main

import (
	"log"

	"github.com/spf13/viper"

	"github.com/MlPablo/CRUDService/internal/server"
)

func main() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	log.Fatal(server.Start())
}
