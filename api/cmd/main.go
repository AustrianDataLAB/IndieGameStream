package main

import (
	"api/controllers"
	"api/repositories"
	"api/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func setupRouter() *gin.Engine {
	//Setup Gin
	r := gin.Default()

	//Setup Repositories
	gamesRepository := repositories.GameRepository()
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

func loadEnv() {
	viper.SetConfigFile("config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	//Load environment file
	loadEnv()

	//Set Gin-gonic to debug or release mode
	gin.SetMode(viper.GetString("GIN_MODE"))

	//Setup Routes
	r := setupRouter()

	// Listen and Server in 0.0.0.0:8080
	r.Run(fmt.Sprintf(":%d", viper.GetInt("port")))
}
