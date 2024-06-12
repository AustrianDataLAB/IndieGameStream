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

func setupRouter(db *sql.DB, azClient *azblob.Client) *gin.Engine {
	//Setup Gin
	r := gin.Default()
	//Cors
	r.Use(CORSMiddleware())

	//Setup Repositories
	gamesRepository := repositories.GameRepository(db)
	gamesService := services.GameService(gamesRepository, azClient)
	gamesController := controllers.GameController(gamesService)

	authService := services.AuthService()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	//Upload a game
	r.POST("/games", authService.Authorize, gamesController.UploadGame)
	//Get all uploaded games
	r.GET("/games", authService.Authorize, gamesController.GetAllGames)
	//Get a specific game by its id
	r.GET("/games/:id", authService.Authorize, gamesController.GetGameById)
	//Delete a specific game, identified by its id
	r.DELETE("/games/:id", authService.Authorize, gamesController.DeleteGameById)

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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func setupAzureBlobContainer(azClient *azblob.Client) {

	containerName := os.Getenv("AZURE_CONTAINER_NAME")
	containerClient := azClient.ServiceClient().NewContainerClient(containerName)

	if containerClient != nil {
		log.Println(fmt.Sprintf("Azure blob container with name %s exists already", containerName))
	} else {
		_, err := azClient.CreateContainer(context.Background(), containerName, nil)
		if err != nil {
			log.Fatal("Creating Azure Blob Container failed")
		}
		log.Println(fmt.Sprintf("Created Azure container %s", containerName))
	}
}

func main() {
	//Load config file
	loadConfig()

	//Setup Azure
	azClient := setupAzureBlobClient()
	setupAzureBlobContainer(azClient)

	//Setup database
	db := setupDatabase()
	defer db.Close()

	//Set Gin-gonic to debug or release mode
	gin.SetMode(os.Getenv("GIN_MODE"))

	//Setup Routes
	r := setupRouter(db, azClient)

	// Listen and Server in 0.0.0.0:8080
	err := r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal(err.Error())
	}
}

func setupAzureBlobClient() *azblob.Client {
	url := fmt.Sprintf("https://%s.blob.core.windows.net/", os.Getenv("AZURE_STORAGE_ACCOUNT"))

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal("Initializing DefaultAzureCredential failed")
	}

	client, err := azblob.NewClient(url, credential, nil)
	if err != nil {
		log.Fatal("Initializing Azure Blob Client failed")
	}

	return client
}
