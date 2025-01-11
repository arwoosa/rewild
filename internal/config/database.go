package config

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDatabase() *mongo.Database {
	dbName := APP.DbApiDatabase
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(APP.DbConnection))

	if err != nil {
		fmt.Println("ERROR", err)
	}

	DB = client.Database(dbName)
	return DB
}
