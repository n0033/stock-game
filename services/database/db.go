package database

import (
	"context"
	"log"
	"time"

	config "github.com/n0033/stock-game/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client, err = mongo.NewClient(options.Client().ApplyURI(config.MONGODB_CONSTR))

var database *mongo.Database

func init() { // this is what makes database variable a singleton
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	database = client.Database(config.MONGODB_DATABASE)
}

func GetDatabase() *mongo.Database {
	return database
}
