package database

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/sankeerthanak/airbnbreplica/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoStorage() (*mongo.Database, error) {

	var client *mongo.Client
	connectionString := "mongodb://" + config.Envs.DBUser + ":" + config.Envs.DBPassword + "@ac-2hiwt8t-shard-00-00.daspgu0.mongodb.net:27017,ac-2hiwt8t-shard-00-01.daspgu0.mongodb.net:27017,ac-2hiwt8t-shard-00-02.daspgu0.mongodb.net:27017/?ssl=true&replicaSet=atlas-sua3z7-shard-0&authSource=admin&retryWrites=true&w=majority&appName=Cluster0"
	clientOption := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB connection successful")

	database := client.Database(config.Envs.DBName)
	//fmt.Println("Here")
	return database, nil
}

func NewRedisStorage() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return client
}
