package database

import (
	"context"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	clientOnce sync.Once
)

// ConnectDB connects to MongoDB
func ConnectDB() {
	mongoURI := os.Getenv("MONGO_URI")
	clientOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(mongoURI)

		// Connect to MongoDB
		var err error
		client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Connected to MongoDB!")
	})
}

// GetCollection returns a collection instance
func GetCollection() *mongo.Collection {
	dbName := os.Getenv("DB_NAME")
	if client == nil {
		log.Fatal("MongoDB client is not initialized")
	}
	return client.Database(dbName).Collection("products")
}
