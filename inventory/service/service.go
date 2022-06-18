package service

import (
	"context"
	"log"
	"net"

	"github.com/joesjo/grpc-store/inventory/database"
	pb "github.com/joesjo/grpc-store/inventory/protobuf"
	"google.golang.org/grpc"
)

const (
	port = ":8080"
)

type server struct {
	pb.UnimplementedInventoryServiceServer
}

type InvalidRequestError struct {
	message string
}

func (e *InvalidRequestError) Error() string {
	return "Invalid request: " + e.message
}

func (s *server) GetInventory(req *pb.Empty, stream pb.InventoryService_GetInventoryServer) error {
	items, err := database.GetAllItems()
	if err != nil {
		return err
	}
	for _, item := range items {
		if err := stream.Send(item); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	if req.Id == "" {
		return nil, &InvalidRequestError{message: "item id is required"}
	}
	item, err := database.FindById(req.Id)
	if err != nil {
		return nil, err
	}
	response := &pb.GetItemResponse{Item: item}
	return response, nil
}

func (s *server) FindItems(req *pb.FindItemsRequest, stream pb.InventoryService_FindItemsServer) error {
	if req.Name == "" {
		return &InvalidRequestError{message: "item name is required"}
	}
	items, err := database.FindByName(req.Name)
	if err != nil {
		return err
	}
	for _, item := range items {
		if err := stream.Send(item); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) InsertItem(ctx context.Context, req *pb.InsertItemRequest) (*pb.InsertItemResponse, error) {
	item := req.Item
	id, err := database.InsertItem(item)
	if err != nil {
		return nil, err
	}
	return &pb.InsertItemResponse{ItemId: id.Hex()}, nil
}

func (s *server) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {
	item := req.Item
	count, err := database.UpdateItem(item)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateItemResponse{Count: count}, nil
}

func (s *server) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	count, err := database.DeleteItem(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteItemResponse{Count: count}, nil
}

func (s *server) IncrementItemQuantity(ctx context.Context, req *pb.IncrementItemQuantityRequest) (*pb.IncrementItemQuantityResponse, error) {
	count, err := database.IncrementItemQuantity(req.Id, req.Amount)
	if err != nil {
		return nil, err
	}
	return &pb.IncrementItemQuantityResponse{Count: count}, nil
}

func Start() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterInventoryServiceServer(s, &server{})
	log.Printf("Starting inventory management server on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
