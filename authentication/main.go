package main

import (
	"github.com/joesjo/grpc-store/authentication/database"
	"github.com/joesjo/grpc-store/authentication/service"
)

func main() {
	database.Init()
	service.Start()
}
