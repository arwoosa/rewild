package main

import (
	"fmt"
	"oosa_rewild/internal/config"
	"oosa_rewild/routes"

	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database

// @title		OOSA API
// @version 	1.0
// @description This is OOSA's API. Generate the swagger documentation by running `swag init --parseDependency`
// @securityDefinitions.apikey 		ApiKeyAuth
// @in								header
// @name 							Authorization
func main() {
	config.InitialiseConfig()
	db = config.ConnectDatabase()

	appPort := config.APP.AppPort
	fmt.Println("Starting app on port: ", appPort)

	engine := routes.RegisterRoutes()
	engine.Run(":" + appPort)
}
