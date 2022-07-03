package serviceclient

import (
	"context"
	"io"

	inventorypb "github.com/joesjo/grpc-store/inventory/protobuf"
	"google.golang.org/grpc"
)

var (
	inventoryClient inventorypb.InventoryServiceClient
)

func Init() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	inventoryClient = inventorypb.NewInventoryServiceClient(conn)
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
