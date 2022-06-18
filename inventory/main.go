package main

import (
	"github.com/joesjo/grpc-store/inventory/database"
	"github.com/joesjo/grpc-store/inventory/service"
)

func main() {
	database.Init()
	service.Start()
}
