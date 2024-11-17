package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ecabigting/letsgo-brrr/usermanager-api/config"
	"github.com/ecabigting/letsgo-brrr/usermanager-api/routes"
	"github.com/ecabigting/letsgo-brrr/usermanager-api/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// load configuration
	config.LoadConfig()
	fmt.Println("Config loaded...")
	// init logger
	utils.InitLogger()
	fmt.Println("Logger Loaded...")
	// connect to db
	clientOptions := options.Client().ApplyURI(config.AppConfig.MongoURI)
	fmt.Println("Loaded DB Client Options...")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println("Error connecting to DB..", err)
		utils.Logger.Fatal("Failed to connect to MongoDB:", err)
		os.Exit(1)
	}
	fmt.Println("DB Connection no error")
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			utils.Logger.Fatal("Failed to disconnect from MongoDB:", err)
		}
	}()

	// setup gin router
	router := gin.Default()

	// setup routes
	routes.SetupRoutes(router, client.Database("usermanager-api"))

	// start server
	if err := router.Run(":" + config.AppConfig.Port); err != nil {
		utils.Logger.Fatal("Fatal error starting the server:", err)
	}
}
