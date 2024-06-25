package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sankeerthanak/airbnbreplica/cmd/api"
	database "github.com/sankeerthanak/airbnbreplica/dataBase"
)

func main() {

	dataBase, err := database.NewMongoStorage()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB connection successful")

	rDatabase := database.NewRedisStorage()
	response, err := rDatabase.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("redis connection failed %s", response)
		log.Fatal(err)
	}

	fmt.Println("Redis connection successful")

	server := api.NewApiServer(":8081", dataBase, rDatabase)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
