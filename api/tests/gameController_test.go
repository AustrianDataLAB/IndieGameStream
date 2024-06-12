package tests

import (
	"api/controllers"
	"api/dtos"
	"api/models"
	"api/repositories"
	"api/services"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"regexp"
	"testing"
)

func Test_Read_Should_Succeed(t *testing.T) {
	//==============Prepare=======================
	db, dbMock := databaseMock()
	defer db.Close()
	prepare_sql_mock_for_Test_Read_Should_Succeed(dbMock)
	gameController := gameController(db)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("subject", "MockOwner")

	//Execute
	gameController.GetAllGames(c)

	//Verify
	if w.Code != 200 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}

	games := []dtos.GetAllGamesResponseBody{}

	err := json.Unmarshal(w.Body.Bytes(), &games)
	if err != nil {
		t.Error(err)
	}
	if len(games) != 2 {
		t.Error(fmt.Sprint("Expected 2 games, got ", len(games)))
	}

	if games[0].ID == uuid.Nil {
		t.Error(fmt.Sprint("Expected non-nil game ID, got ", games[0].ID))
	}
	if games[1].ID == uuid.Nil {
		t.Error(fmt.Sprint("Expected non-nil game ID, got ", games[1].ID))
	}
	if len(games[0].Url) == 0 {
		t.Error(fmt.Sprint("Expected non-empty game URL, got ", games[0].Url))
	}
	if len(games[1].Url) == 0 {
		t.Error(fmt.Sprint("Expected non-empty game URL, got ", games[1].Url))
	}
	if len(games[0].Title) == 0 {
		t.Error(fmt.Sprint("Expected non-empty game title, got ", games[0].Title))
	}
	if len(games[1].Title) == 0 {
		t.Error(fmt.Sprint("Expected non-empty game title, got ", games[1].Title))
	}
}

func prepare_sql_mock_for_Test_Read_Should_Succeed(mock sqlmock.Sqlmock) {
	gameA := models.Game{
		ID:              uuid.New(),
		Title:           "AMockA",
		StorageLocation: "AMockB",
		Status:          "AMockC",
		Url:             "AMockD",
		Owner:           "MockOwner",
	}

	gameB := models.Game{
		ID:              uuid.New(),
		Title:           "BMockA",
		StorageLocation: "BMockB",
		Status:          "BMockC",
		Url:             "BMockD",
		Owner:           "MockOwner",
	}

	mock.ExpectPrepare(regexp.QuoteMeta("SELECT * FROM games WHERE owner = ?"))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM games WHERE owner = ?")).
		WithArgs("MockOwner").
		WillReturnRows(
			sqlmock.NewRows([]string{"ID", "Title", "StorageLocation", "Status", "Url", "Owner"}).
				AddRow(gameA.ID, gameA.Title, gameA.StorageLocation, gameA.Status, gameA.Url, gameA.Owner).
				AddRow(gameB.ID, gameB.Title, gameB.StorageLocation, gameB.Status, gameB.Url, gameB.Owner),
		)
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
	gamesService := services.GameService(gamesRepository)
	return controllers.GameController(gamesService)
}
