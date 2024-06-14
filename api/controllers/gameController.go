package controllers

import (
	"api/dtos"
	"api/services"
	"database/sql"
	"fmt"
	"github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
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
	games, err := g.service.FindAllByOwner(c.GetString("subject"))
	if err != nil { //TODO handle different errors
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	//Map to dto
	resultDto := []dtos.GetGameByIdResponseBody{}
	err = dto.Map(&resultDto, games)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, resultDto)
	return
}

func (g gameController) GetGameById(c *gin.Context) {
	_uuid := getUUIDFromRequest(c)
	if _uuid != uuid.Nil {
		//Get game by uuid
		game, err := g.service.FindByID(_uuid)
		if err != nil { //TODO handle different errors
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if game == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Game not found"})
			return
		}
		if game.Owner != c.GetString("subject") {
			log.Print(fmt.Printf("%s tried to access an resource of %s", c.GetString("subject"), game.Owner))
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"message": "You don't have permission to access this resource"})
			return
		}

		//Map to dto
		resultDto := dtos.GetGameByIdResponseBody{}
		err = dto.Map(&resultDto, game)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, resultDto)
		return
	}
}

func (g gameController) UploadGame(c *gin.Context) {

	//Try to read the title from body
	title := c.Request.PostFormValue("title")
	if len(title) == 0 {
		//If it is not in the body, check if it is a query parameter
		title = c.GetString("title")
		//If the title is still empty return BadRequest
		if len(title) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Title is required"})
			return
		}
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	sub := c.GetString("subject")
	if len(sub) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "IdToken is invalid, sub is missing"})
		return
	}

	//Save the game in the database and azure
	game, err := g.service.Save(file, title, sub)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Header("content-location", fmt.Sprintf("%s/games/%s", c.Request.Host, game.ID.String()))
	c.AbortWithStatus(http.StatusCreated)
	return
}

func (g gameController) DeleteGameById(c *gin.Context) {
	_uuid := getUUIDFromRequest(c)
	if _uuid != uuid.Nil {

		//Check if the user has access to the game
		authorized, err := g.hasAccessToGame(c)

		if err != nil {
			if err == sql.ErrNoRows {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Game not found"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		if !authorized {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "You don't have permission to access this resource"})
			return
		}

		err = g.service.Delete(_uuid)
		if err != nil { //TODO handle different errors
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
	}
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid game ID"})
		return uuid.Nil
	} else if _uuid == uuid.Nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid game ID"})
		return uuid.Nil
	}
	return _uuid
}

// Returns true if the owner the user who is logged-in has the same subject-id as the game owner.
// Returns false and error if any other error occurred.
func (g gameController) hasAccessToGame(c *gin.Context) (bool, error) {
	owner, err := g.service.ReadOwner(getUUIDFromRequest(c))
	if err != nil {
		return false, err
	}

	return owner == c.GetString("subject"), nil
}
