package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
		log.Fatal(err)
	}

	const seconds time.Duration = 10 // seconds
	ctx, cancel := context.WithTimeout(context.Background(), seconds*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Panic(err)
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Connected to MongoDB")

	// Create indices
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	indexName, err := GetCollection(client, "users").Indexes().CreateOne(
		context.TODO(),
		indexModel,
	)

	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Name of Index Created: " + indexName)

	return client
}

// DB Client instance
var DB *mongo.Client = ConnectDB()

// Getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("twitter").Collection(collectionName)
	return collection
}
