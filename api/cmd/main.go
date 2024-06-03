package main

import (
	"api/controllers"
	"api/repositories"
	"api/scripts"
	"api/services"
	"context"
	"database/sql"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func setupRouter(db *sql.DB) *gin.Engine {
	//Setup Gin
	r := gin.Default()

	//Setup Repositories
	gamesRepository := repositories.GameRepository(db)
	gamesService := services.GameService(gamesRepository)
	gamesController := controllers.GameController(gamesService)

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	//Upload a game
	r.POST("/games/", gamesController.UploadGame)
	//Get all uploaded games
	r.GET("/games", gamesController.GetAllGames)
	//Get a specific game by its id
	r.GET("/games/:id", gamesController.GetGameById)
	//Delete a specific game, identified by its id
	r.DELETE("/games/:id", gamesController.DeleteGameById)

	return r
}

func loadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(fmt.Sprintf("Failed to load .env file: %s", err))
		log.Println("Environment variables will be used")
	}
}

func setupDatabase() *sql.DB {
	//Create database if it is not existing yet.
	//We might have to remove this if we use an azure database
	scripts.CreateDatabaseIfNotExists(os.Getenv("MYSQL_DATABASE"))
	//Connect to the database
	db := scripts.ConnectToDatabase()
	//Check if database is online
	err := db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
	//Check if we have new migrations and apply them
	scripts.MigrateDatabase(db)
	return db
}

func setupAzure() {
	client, err := azureClient()
	if err != nil {
		log.Fatal(err)
	}

	//TODO: Check if container is already existing

	containerName := os.Getenv("AZURE_CONTAINER_NAME")
	_, err = client.CreateContainer(context.Background(), containerName, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("Created Azure container %s", containerName))

}

func main() {
	//Load config file
	loadConfig()

	//Setup Azure
	setupAzure()

	//Setup database
	db := setupDatabase()
	defer db.Close()

	//Set Gin-gonic to debug or release mode
	gin.SetMode(os.Getenv("GIN_MODE"))

	//Setup Routes
	r := setupRouter(db)

	// Listen and Server in 0.0.0.0:8080
	err := r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal(err.Error())
	}
}

func azureClient() (*azblob.Client, error) {
	url := fmt.Sprintf("https://%s.blob.core.windows.net/", os.Getenv("AZURE_STORAGE_ACCOUNT"))

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	return azblob.NewClient(url, credential, nil)
}
