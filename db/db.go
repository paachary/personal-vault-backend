package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoClient *mongo.Client
var DatabaseName string
var CollectionName string
var MFACollectionName string

func Init() {

	mongoHost := os.Getenv("MONGO_DB_HOST")
	mongoPort := os.Getenv("MONGO_DB_PORT")

	uri := fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort)

	var err error

	bsonOpts := options.BSONOptions{NilSliceAsEmpty: true}
	clientOptions := options.Client().ApplyURI(uri).SetBSONOptions(&bsonOpts)

	MongoClient, err = mongo.Connect(clientOptions)

	DatabaseName = os.Getenv("MONGO_DB_NAME")

	CollectionName = os.Getenv("MONGO_COLLECTION_NAME")

	MFACollectionName = os.Getenv("MONGO_MFA_COLLECTION_NAME")

	if err != nil {
		panic(err)
	}

	err = MongoClient.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	// Successfully connected and pinged the MongoDB server.
	fmt.Println("Connected to MongoDB!")

	err = createCollections()

	if err != nil {
		panic(err)
	}

}
