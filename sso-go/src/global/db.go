package global

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB mongo.Database
var Ctx context.Context

func ConnectToMongo() {

	client, err := mongo.NewClient(options.Client().ApplyURI(DB_URL))
	if err != nil {
		log.Println("mongodb error")
		log.Printf("error: %v", err)
	}
	Ctx = context.Background()

	ctx, cancel := context.WithTimeout(Ctx, 10*time.Second)
	defer cancel()

	client.Connect(ctx)

	DB = *client.Database("shopping")
}
