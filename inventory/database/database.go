package database

import (
	"context"
	"log"
	"time"

	"github.com/jinzhu/copier"
	pb "github.com/joesjo/grpc-store/inventory/protobuf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	mongouri       = "mongodb://localhost:2717"
	databaseName   = "store"
	collectionName = "inventory"
)

var (
	client     *mongo.Client
	collection *mongo.Collection
	err        error
)

type InventoryItem struct {
	Id               primitive.ObjectID `bson:"_id,omitempty"`
	pb.InventoryItem `bson:",inline"`
}

func Init() {
	client, err = mongo.NewClient(options.Client().ApplyURI(mongouri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not connect to mongodb server on: ", mongouri)
	}
	collection = client.Database(databaseName).Collection(collectionName)
}

func find(filter primitive.D) ([]*pb.InventoryItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var result []*pb.InventoryItem
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var (
			item       InventoryItem
			resultItem = &pb.InventoryItem{}
		)
		err := cursor.Decode(&item)
		if err != nil {
			return nil, err
		}
		copier.Copy(resultItem, &item.InventoryItem)
		resultItem.Id = item.Id.Hex()
		result = append(result, resultItem)
	}
	return result, nil
}

func GetAllItems() ([]*pb.InventoryItem, error) {
	log.Println("Getting all items")
	filter := bson.D{}
	return find(filter)
}

func InsertItem(item *pb.InventoryItem) (primitive.ObjectID, error) {
	log.Println("Inserting item:", item)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var newItem = &InventoryItem{}
	copier.Copy(newItem, item)
	result, err := collection.InsertOne(ctx, newItem)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	resultId := result.InsertedID.(primitive.ObjectID)
	return resultId, nil
}

func FindById(itemId string) (*pb.InventoryItem, error) {
	log.Println("Finding item by id:", itemId)
	objId, err := primitive.ObjectIDFromHex(itemId)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: objId}}
	res, err := find(filter)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, status.Errorf(codes.NotFound, "Could not find item with id "+itemId)
	}
	return res[0], nil
}

func FindByName(name string) ([]*pb.InventoryItem, error) {
	log.Println("Finding item by name:", name)
	filter := bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: name}}}}
	return find(filter)
}

func UpdateItem(item *pb.InventoryItem) (int64, error) {
	log.Println("Updating item:", item)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var newItem = &InventoryItem{}
	copier.Copy(newItem, item)
	update := bson.D{
		{Key: "$set", Value: newItem},
	}
	objId, err := primitive.ObjectIDFromHex(item.Id)
	res, err := collection.UpdateByID(ctx, objId, update)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, err
}

func DeleteItem(itemId string) (int64, error) {
	log.Println("Deleting item:", itemId)
	objId, err := primitive.ObjectIDFromHex(itemId)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "_id", Value: objId}}
	res, err := collection.DeleteOne(ctx, filter)
	return res.DeletedCount, err
}

func IncrementItemQuantity(itemId string, quantity int32) (int64, error) {
	log.Println("Incrementing item quantity:", itemId, quantity)
	objId, err := primitive.ObjectIDFromHex(itemId)
	if err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "_id", Value: objId}}
	update := bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "quantity", Value: quantity},
		}},
	}
	res, err := collection.UpdateOne(ctx, filter, update)
	return res.ModifiedCount, err
}
