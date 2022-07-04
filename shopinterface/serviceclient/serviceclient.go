package serviceclient

import (
	"context"
	"io"
	"log"
	"os"

	authenticationpb "github.com/joesjo/grpc-store/authentication/protobuf"
	inventorypb "github.com/joesjo/grpc-store/inventory/protobuf"
	"google.golang.org/grpc"
)

const (
	INVENTORY_URI      = "localhost:8082"
	AUTHENTICATION_URI = "localhost:8081"
)

var (
	inventoryClient      inventorypb.InventoryServiceClient
	authenticationClient authenticationpb.AuthenticationServiceClient
)

func Init() {
	go func() {
		inventoryurl, exists := os.LookupEnv("INVENTORY_URI")
		if !exists {
			inventoryurl = INVENTORY_URI
		}
		inventoryClient = connectInventory(inventoryurl)
	}()
	go func() {
		authenticationurl, exists := os.LookupEnv("AUTHENTICATION_URI")
		if !exists {
			authenticationurl = AUTHENTICATION_URI
		}
		authenticationClient = connectAuthentication(authenticationurl)
	}()
}

func connectInventory(url string) inventorypb.InventoryServiceClient {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	log.Println("Connected to inventory service")
	return inventorypb.NewInventoryServiceClient(conn)
}

func connectAuthentication(url string) authenticationpb.AuthenticationServiceClient {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	log.Println("Connected to authentication service")
	return authenticationpb.NewAuthenticationServiceClient(conn)
}

// Get inventory using grpc stream
func GetInventory() ([]*inventorypb.InventoryItem, error) {
	stream, err := inventoryClient.GetInventory(context.Background(), &inventorypb.Empty{})
	if err != nil {
		return nil, err
	}
	var items []*inventorypb.InventoryItem
	for {
		item, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func GetItem(itemId string) (*inventorypb.InventoryItem, error) {
	itemRequest := &inventorypb.GetItemRequest{Id: itemId}
	item, err := inventoryClient.GetItem(context.Background(), itemRequest)
	return item.GetItem(), err
}

func FindItems(name string) ([]*inventorypb.InventoryItem, error) {
	itemRequest := &inventorypb.FindItemsRequest{Name: name}
	stream, err := inventoryClient.FindItems(context.Background(), itemRequest)
	if err != nil {
		return nil, err
	}
	var items []*inventorypb.InventoryItem
	for {
		item, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func CreateItem(name string, quantity int32) (string, error) {
	itemRequest := &inventorypb.InsertItemRequest{Item: &inventorypb.InventoryItem{Name: name, Quantity: quantity}}
	itemId, err := inventoryClient.InsertItem(context.Background(), itemRequest)
	return itemId.GetItemId(), err
}

func StockItem(itemId string, quantity int32) error {
	itemRequest := &inventorypb.IncrementItemQuantityRequest{Id: itemId, Amount: quantity}
	_, err := inventoryClient.IncrementItemQuantity(context.Background(), itemRequest)
	return err
}

func PurchaseItem(itemId string, quantity int32) error {
	itemRequest := &inventorypb.IncrementItemQuantityRequest{Id: itemId, Amount: -quantity}
	_, err := inventoryClient.IncrementItemQuantity(context.Background(), itemRequest)
	return err
}

func UpdateItem(itemId string, name string, quantity int32) error {
	itemRequest := &inventorypb.UpdateItemRequest{Item: &inventorypb.InventoryItem{Id: itemId, Name: name, Quantity: quantity}}
	_, err := inventoryClient.UpdateItem(context.Background(), itemRequest)
	return err
}

func DeleteItem(itemId string) error {
	itemRequest := &inventorypb.DeleteItemRequest{Id: itemId}
	_, err := inventoryClient.DeleteItem(context.Background(), itemRequest)
	return err
}

func CreateUser(username string, password string) (string, error) {
	userRequest := &authenticationpb.CreateUserRequest{User: &authenticationpb.User{Username: username, Password: password}}
	userId, err := authenticationClient.CreateUser(context.Background(), userRequest)
	return userId.GetError(), err
}

func Login(username string, password string) (string, error) {
	userRequest := &authenticationpb.AuthenticateRequest{User: &authenticationpb.User{Username: username, Password: password}}
	userId, err := authenticationClient.Authenticate(context.Background(), userRequest)
	return userId.GetToken(), err
}

func ValidateToken(token string) (string, error) {
	userRequest := &authenticationpb.ValidateTokenRequest{Token: token}
	userId, err := authenticationClient.ValidateToken(context.Background(), userRequest)
	return userId.GetUsername(), err
}
