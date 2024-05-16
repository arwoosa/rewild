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
	dbHost := APP.DbApiHost
	dbPort := APP.DbApiPort
	dbName := APP.DbApiDatabase
	//dbUser := APP.DbApiUsername
	//dbPassword := APP.DbApiPassword
	//connectionString := "mongodb://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/?timeoutMS=5000&maxPoolSize=20&w=majority"
	connectionString := "mongodb://" + dbHost + ":" + dbPort + "/?timeoutMS=5000&maxPoolSize=20&w=majority"
	fmt.Println(dbHost, dbName, connectionString)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))

	if err != nil {
		fmt.Println("ERROR", err)
	}

	DB = client.Database(dbName)
	return client.Database(dbName)
}
