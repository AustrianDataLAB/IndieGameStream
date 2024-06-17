package tests

import (
	"api/controllers"
	"api/dtos"
	"api/models"
	"api/repositories"
	"api/services"
	"api/tests/mocks"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"regexp"
	"testing"
)

func Test_Read_By_Id_Should_Succeed(t *testing.T) {
	//======================= PREPARE	PREPARE		PREPARE		PREPARE =======================
	owner := "MockOwner"
	// Create database mock
	db, dbMock := databaseMock()
	defer db.Close()
	// Create Models
	game := mocks.GameMock("A")
	game.Owner = owner
	// Define queries
	dbMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM games WHERE ID = ?")).
		WithArgs(game.ID).WillReturnRows(
		sqlmock.NewRows([]string{"Id", "Title", "StorageLocation", "Status", "Url", "Owner"}).
			AddRow(game.ID, game.Title, game.StorageLocation, game.Status, game.Url, game.Owner),
	)

	// Finally, create gameController
	gameController := gameController(db)
	// Prepare Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("subject", owner)

	//======================= EXECUTE	EXECUTE		EXECUTE		EXECUTE =======================
	c.Params = gin.Params{gin.Param{Key: "id", Value: game.ID.String()}}
	gameController.GetGameById(c)

	//======================= VERIFY	VERIFY		VERIFY		VERIFY =======================
	//Check HTTP response
	if w.Code != 200 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}

	//Check response body
	var responseBody dtos.GetAllGamesResponseBody
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	if &responseBody == nil {
		t.Error("response body is empty")
	}

	//Verify games
	verifyDto(t, &responseBody, game)

}

func Test_Read_All_Should_Succeed(t *testing.T) {
	//======================= PREPARE	PREPARE		PREPARE		PREPARE =======================
	owner := "MockOwner"
	// Create database mock
	db, dbMock := databaseMock()
	defer db.Close()
	// Create Models
	gameA := mocks.GameMock("A")
	gameA.Owner = owner
	gameB := mocks.GameMock("B")
	gameB.Owner = owner
	// Define queries
	dbMock.ExpectPrepare(regexp.QuoteMeta("SELECT * FROM games WHERE owner = ?"))
	dbMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM games WHERE owner = ?")).
		WithArgs(owner).
		WillReturnRows(
			sqlmock.NewRows([]string{"ID", "Title", "StorageLocation", "Status", "Url", "Owner"}).
				AddRow(gameA.ID, gameA.Title, gameA.StorageLocation, gameA.Status, gameA.Url, gameA.Owner).
				AddRow(gameB.ID, gameB.Title, gameB.StorageLocation, gameB.Status, gameB.Url, gameB.Owner),
		)
	// Finally, create gameController
	gameController := gameController(db)
	// Prepare Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("subject", owner)

	//======================= EXECUTE	EXECUTE		EXECUTE		EXECUTE =======================
	gameController.GetAllGames(c)

	//======================= VERIFY	VERIFY		VERIFY		VERIFY =======================
	//Check HTTP response
	if w.Code != 200 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}

	//Check response body
	var responseBody []dtos.GetAllGamesResponseBody
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	if len(responseBody) != 2 {
		t.Error(fmt.Sprint("Expected 2 games, got ", len(responseBody)))
	}

	//Verify games
	verifyDto(t, &responseBody[0], gameA)
	verifyDto(t, &responseBody[1], gameB)

}

func verifyDto(t *testing.T, dto *dtos.GetAllGamesResponseBody, game *models.Game) {
	if dto.Title != game.Title {
		t.Error(fmt.Sprintf("Expected title %s, got %s", game.Title, dto.Title))
	}
	if dto.Url != game.Url {
		t.Error(fmt.Sprintf("Expected url %s, got %s", game.Url, dto.Url))
	}
	if dto.ID != game.ID {
		t.Error(fmt.Sprintf("Expected id %v, got %v", game.ID, dto.ID))
	}
	if dto.Status != game.Status {
		t.Error(fmt.Sprintf("Expected status %v, got %v", game.Status, dto.Status))
	}
}

func databaseMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return db, mock
}

func gameController(db *sql.DB) controllers.IGameController {
	gamesRepository := repositories.GameRepository(db)
	gamesService := services.GameService(gamesRepository, nil, nil)
	return controllers.GameController(gamesService)
}
