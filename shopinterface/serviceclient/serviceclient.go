package serviceclient

import (
	inventorypb "github.com/joesjo/grpc-store/inventory/protobuf"
	"google.golang.org/grpc"
)

var (
	inventoryClient inventorypb.InventoryServiceClient
)

func Init() {
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	inventoryClient = inventorypb.NewInventoryServiceClient(conn)
}
