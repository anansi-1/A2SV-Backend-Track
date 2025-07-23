package main

import (
	"context"
	"log"
	router "task-seven/Delivery/routers"

	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	clientOption := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("taskdb")
	router.NewTaskRouter(10*time.Second, *db)
}
