package main

import (
	"log"

	"github.com/MlPablo/CRUDService/internal/server"
)

func main() {
	log.Fatal(server.Start())
}
