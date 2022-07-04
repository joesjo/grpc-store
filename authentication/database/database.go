package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongouri       = "mongodb://localhost:2717"
	databaseName   = "store"
	collectionName = "users"
)

var (
	client     *mongo.Client
	collection *mongo.Collection
	err        error
)

type User struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
}

func Init() {
	url, exists := os.LookupEnv("MONGO_URI")
	if !exists {
		url = mongouri
	}
	client, err = mongo.NewClient(options.Client().ApplyURI(url))
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

func CreateUser(username string, password string) error {
	user := User{
		Username: username,
		Password: password,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func FindUser(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user User
	err = collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
