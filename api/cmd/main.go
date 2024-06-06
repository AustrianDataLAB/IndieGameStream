package main

import (
	"api/controllers"
	"api/repositories"
	"api/scripts"
	"api/services"
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
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

	//Cors
	r.Use(cors.Default())

	//Setup Repositories
	gamesRepository := repositories.GameRepository(db)
	gamesService := services.GameService(gamesRepository)
	gamesController := controllers.GameController(gamesService)

	authService := services.AuthService()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	//Upload a game
	r.POST("/games", CorsHeader, authService.Authorize, gamesController.UploadGame)
	//Get all uploaded games
	r.GET("/games", CorsHeader, authService.Authorize, gamesController.GetAllGames)
	//Get a specific game by its id
	r.GET("/games/:id", CorsHeader, authService.Authorize, gamesController.GetGameById)
	//Delete a specific game, identified by its id
	r.DELETE("/games/:id", CorsHeader, authService.Authorize, gamesController.DeleteGameById)

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

func main() {
	//Load config file
	loadConfig()

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

func CorsHeader(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
