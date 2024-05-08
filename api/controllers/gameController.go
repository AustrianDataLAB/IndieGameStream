package controllers

import (
	"api/dtos"
	"api/services"
	"database/sql"
	"fmt"
	"github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type IGameController interface {
	GetAllGames(c *gin.Context)
	GetGameById(c *gin.Context)
	UploadGame(c *gin.Context)
	DeleteGameById(c *gin.Context)
}

type gameController struct {
	service services.IGameService
}

func (g gameController) GetAllGames(c *gin.Context) {
	//Get Games
	games, err := g.service.FindAll()
	if err != nil { //TODO handle different errors
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	//Map to dto
	resultDto := []dtos.GetGameByIdResponseBody{}
	err = dto.Map(&resultDto, games)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, resultDto)
}

func (g gameController) GetGameById(c *gin.Context) {
	_uuid := getUUIDFromRequest(c)
	if _uuid != uuid.Nil {
		//Get game by uuid
		game, err := g.service.FindByID(_uuid)
		if err != nil { //TODO handle different errors
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if game == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Game not found"})
			return
		}

		//Map to dto
		resultDto := dtos.GetGameByIdResponseBody{}
		err = dto.Map(&resultDto, game)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, resultDto)
	}
}

func (g gameController) UploadGame(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	_, err = g.service.Save(file)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.Header("content-location", fmt.Sprintf("%s/games/%s", c.Request.Host, file.Filename))
	c.String(http.StatusCreated, "")
}

func (g gameController) DeleteGameById(c *gin.Context) {
	_uuid := getUUIDFromRequest(c)
	if _uuid != uuid.Nil {
		err := g.service.Delete(_uuid)
		if err != nil { //TODO handle different errors
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"message": "Game not found"})
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		} else {
			c.Status(http.StatusNoContent)
		}
	}

	//TODO Implement delete game
	c.String(http.StatusNoContent, "")
}

func GameController(service services.IGameService) IGameController {
	return &gameController{
		service: service,
	}
}

// Parses the UUID from the request param "uuid" and returns it.
// It returns HTTP 400 and uuid.nil if the uuid is invalid or null
func getUUIDFromRequest(c *gin.Context) uuid.UUID {
	_uuid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid game ID"})
		return uuid.Nil
	} else if _uuid == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid game ID"})
		return uuid.Nil
	}
	return _uuid
}
