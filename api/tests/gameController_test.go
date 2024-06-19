package tests

import (
	"api/apis"
	"api/controllers"
	"api/dtos"
	"api/models"
	"api/repositories"
	"api/services"
	"api/shared"
	"api/tests/mocks"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	v1 "indiegamestream.com/indiegamestream/api/stream/v1"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"regexp"
	"testing"
)

func Test_Read_By_Id_And_Refresh_Should_Succeed(t *testing.T) {
	//======================= PREPARE	PREPARE		PREPARE		PREPARE =======================
	owner := "MockOwner"
	url := "fsdfsf-91f3975e-cdd5-4b9b-8f00-a86bb44d4b82.possum-climb.ts.net"
	// Create Models
	game := mocks.GameMock("A")
	game.ID, _ = uuid.Parse("66c887ca-1f56-426e-ac0c-bc92fff8b798")
	game.Owner = owner
	game.Url = ""
	game.Status = shared.Status_Installing
	// Create fake k8s client
	fakek8s := mocks.K8sMock(&mock.Mock{})
	k8sApi := apis.K8sService(fakek8s)
	//Mock the calls to k8s
	fakek8s.Mock().
		On("Get",
			mock.AnythingOfType("context.backgroundCtx"),
			mock.AnythingOfType("types.NamespacedName"),
			mock.AnythingOfType("*v1.Game")).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(2).(*v1.Game)
			arg.Status.URL = url
		})

	// Create database mock
	db, dbMock := databaseMock()
	defer db.Close()
	// Define queries
	dbMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM games WHERE ID = ?")).
		WithArgs(game.ID).WillReturnRows(
		sqlmock.NewRows([]string{"Id", "Title", "StorageLocation", "Status", "Url", "Owner", "FileName"}).
			AddRow(game.ID, game.Title, game.StorageLocation, game.Status, game.Url, game.Owner, game.FileName),
	)

	dbMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM games WHERE ID = ?")).
		WithArgs(game.ID).WillReturnRows(
		sqlmock.NewRows([]string{"Id", "Title", "StorageLocation", "Status", "Url", "Owner", "FileName"}).
			AddRow(game.ID, game.Title, game.StorageLocation, game.Status, game.Url, game.Owner, game.FileName),
	)

	dbMock.ExpectPrepare(regexp.
		QuoteMeta("UPDATE games SET Title=?, StorageLocation=?, Status=?, Url=?, FileName=? WHERE ID = ?"))

	dbMock.ExpectExec(regexp.
		QuoteMeta("UPDATE games SET Title=?, StorageLocation=?, Status=?, Url=?, FileName=? WHERE ID = ?")).
		WithArgs(game.Title, game.StorageLocation, shared.Status_Installed, url, game.FileName, game.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Finally, create gameController
	gameController := gameController(db, k8sApi, nil)

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
	game.Url = url                        // Game url should be set correctly
	game.Status = shared.Status_Installed //Game Status should be set correctly
	verifyDto(t, &responseBody, game)

}

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
		sqlmock.NewRows([]string{"Id", "Title", "StorageLocation", "Status", "Url", "Owner", "FileName"}).
			AddRow(game.ID, game.Title, game.StorageLocation, game.Status, game.Url, game.Owner, game.FileName),
	)

	// Finally, create gameController
	gameController := gameController(db, nil, nil)
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
			sqlmock.NewRows([]string{"ID", "Title", "StorageLocation", "Status", "Url", "Owner", "FileName"}).
				AddRow(gameA.ID, gameA.Title, gameA.StorageLocation, gameA.Status, gameA.Url, gameA.Owner, gameA.FileName).
				AddRow(gameB.ID, gameB.Title, gameB.StorageLocation, gameB.Status, gameB.Url, gameB.Owner, gameB.FileName),
		)
	// Finally, create gameController
	gameController := gameController(db, nil, nil)
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

func gameController(db *sql.DB, k8s apis.IK8sApi, azure apis.IAzureApi) controllers.IGameController {
	gamesRepository := repositories.GameRepository(db)
	gamesService := services.GameService(gamesRepository, k8s, azure)
	return controllers.GameController(gamesService)
}
