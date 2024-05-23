package tests

import (
	"api/models"
	"api/repositories"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"log"
	"regexp"
	"testing"
)

// ************************************ BEGIN DELETE TESTS ************************************
func Test_Delete_Should_Succeed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer db.Close()

	//Define the mock
	id := uuid.New()
	mock.ExpectPrepare(regexp.QuoteMeta("DELETE FROM games WHERE ID = ?"))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM games")).
		WithArgs(id.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	//Run the test
	repository := repositories.GameRepository(db)

	err = repository.Delete(id)
	if err != nil {
		t.Errorf(err.Error())
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf(err.Error())
	}
}

func Test_Delete_Not_Existing_Should_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer db.Close()

	//Define the mock
	id := uuid.New()
	mock.ExpectPrepare(regexp.QuoteMeta("DELETE FROM games WHERE ID = ?"))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM games")).
		WithArgs(id.String()).
		WillReturnResult(sqlmock.NewResult(0, 0))

	//Run the test
	repository := repositories.GameRepository(db)

	err = repository.Delete(id)

	if err == nil {
		t.Errorf("error was not returned")
	}
	if err != sql.ErrNoRows {
		t.Errorf("wrong error was returned")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf(err.Error())
	}
}

//************************************ END DELETE TESTS ************************************

// ************************************ BEGIN INSERT TESTS ************************************
func Test_Create_Game_Without_Id_Should_Succeed_And_SetId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer db.Close()

	//Define the mock
	game := models.Game{
		ID:              uuid.Nil,
		Title:           "",
		StorageLocation: "",
		Status:          "",
		Url:             "",
	}
	mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO games"))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO games")).
		WithArgs(sqlmock.AnyArg(), game.Title, game.StorageLocation, game.Status, game.Url).
		WillReturnResult(sqlmock.NewResult(0, 1))

	//Run the test
	repository := repositories.GameRepository(db)

	err = repository.Save(&game)
	if err != nil {
		t.Errorf(err.Error())
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf(err.Error())
	}

	if game.ID == uuid.Nil {
		t.Errorf("game id was not created")
	}

}

func Test_Create_Game_With_Id_Should_Succeed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer db.Close()

	//Define the mock
	id := uuid.New()
	game := models.Game{
		ID:              id,
		Title:           "",
		StorageLocation: "",
		Status:          "",
		Url:             "",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM games WHERE ID = ?")).
		WithArgs(id).WillReturnError(sql.ErrNoRows)

	mock.ExpectPrepare("INSERT INTO games")

	mock.ExpectExec("INSERT INTO games").
		WithArgs(game.ID, game.Title, game.StorageLocation, game.Status, game.Url).
		WillReturnResult(sqlmock.NewResult(0, 1))

	//Run the test
	repository := repositories.GameRepository(db)

	err = repository.Save(&game)
	if err != nil {
		t.Errorf(err.Error())
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf(err.Error())
	}

	if game.ID != id {
		t.Errorf("game id has been changed")
	}

}

//************************************ END INSERT TESTS ************************************

//************************************ BEGIN UPDATE TESTS ************************************

func Test_Save_Existing_Game_Should_Succeed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer db.Close()

	id := uuid.New()
	game := models.Game{
		ID:              id,
		Title:           "Mock",
		StorageLocation: "Mock",
		Status:          "Mock",
		Url:             "Mock",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM games WHERE ID = ?")).
		WithArgs(id).WillReturnRows(
		sqlmock.NewRows([]string{"id", "title", "storage_location", "status", "url"}).
			AddRow(id, "", "", "", ""),
	)

	mock.ExpectPrepare(regexp.QuoteMeta("UPDATE games SET Title=?, StorageLocation=?, Status=?, Url=? WHERE ID = ?"))

	mock.ExpectExec(regexp.QuoteMeta("UPDATE games SET Title=?, StorageLocation=?, Status=?, Url=? WHERE ID = ?")).
		WithArgs(game.Title, game.StorageLocation, game.Status, game.Url, game.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	//Run the test
	repository := repositories.GameRepository(db)

	err = repository.Save(&game)
	if err != nil {
		t.Errorf(err.Error())
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf(err.Error())
	}

	if game.ID != id {
		t.Errorf("game id has been changed")
	}

}

//************************************ END UPDATE TESTS ************************************

//************************************ BEGIN READ TESTS ************************************

func Test_Find_Game_By_Id_Should_Succeed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer db.Close()

	id := uuid.New()
	game := models.Game{
		ID:              id,
		Title:           "MockA",
		StorageLocation: "MockB",
		Status:          "MockC",
		Url:             "MockD",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM games WHERE ID = ?")).
		WithArgs(id).WillReturnRows(
		sqlmock.NewRows([]string{"id", "title", "storage_location", "status", "url"}).
			AddRow(id, game.Title, game.StorageLocation, game.Status, game.Url),
	)

	//Run the test
	repository := repositories.GameRepository(db)

	res, err := repository.FindByID(game.ID)
	if err != nil {
		t.Errorf(err.Error())
	}

	if res == nil {
		t.Errorf("game was not returned")
	}

	if game.ID != res.ID {
		t.Errorf("game id has been changed")
	}

	if game.Title != res.Title {
		t.Errorf("game title has been changed")
	}

	if game.StorageLocation != res.StorageLocation {
		t.Errorf("game storage location has been changed")
	}

	if game.Status != res.Status {
		t.Errorf("game status has been changed")
	}

	if game.Url != res.Url {
		t.Errorf("game url has been changed")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf(err.Error())
	}
}

func Test_Find_Game_By_Id_Should_Return_Nil(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer db.Close()

	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM games WHERE ID = ?")).
		WithArgs(id).WillReturnError(sql.ErrNoRows)

	//Run the test
	repository := repositories.GameRepository(db)

	res, err := repository.FindByID(id)
	if err != nil {
		t.Errorf(err.Error())
	}

	if res != nil {
		t.Errorf("something has been returned, but it should not return anything.")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf(err.Error())
	}
}

func Test_Find_Two_Games_Should_Succeed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer db.Close()

	gameA := models.Game{
		ID:              uuid.New(),
		Title:           "AMockA",
		StorageLocation: "AMockB",
		Status:          "AMockC",
		Url:             "AMockD",
	}

	gameB := models.Game{
		ID:              uuid.New(),
		Title:           "BMockA",
		StorageLocation: "BMockB",
		Status:          "BMockC",
		Url:             "BMockD",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM games")).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "title", "storage_location", "status", "url"}).
				AddRow(gameA.ID, gameA.Title, gameA.StorageLocation, gameA.Status, gameA.Url).
				AddRow(gameB.ID, gameB.Title, gameB.StorageLocation, gameB.Status, gameB.Url),
		)

	//Run the test
	repository := repositories.GameRepository(db)

	res, err := repository.FindAll()
	if err != nil {
		t.Errorf(err.Error())
	}

	if res == nil {
		t.Errorf("game was not returned")
	}

	if len(res) != 2 {
		t.Errorf("FindAll should return two games")
	}

	if gameA.ID != res[0].ID {
		t.Errorf("game id has been changed")
	}

	if gameA.Title != res[0].Title {
		t.Errorf("game title has been changed")
	}

	if gameA.StorageLocation != res[0].StorageLocation {
		t.Errorf("game storage location has been changed")
	}

	if gameA.Status != res[0].Status {
		t.Errorf("game status has been changed")
	}

	if gameA.Url != res[0].Url {
		t.Errorf("game url has been changed")
	}

	if gameB.ID != res[1].ID {
		t.Errorf("game id has been changed")
	}

	if gameB.Title != res[1].Title {
		t.Errorf("game title has been changed")
	}

	if gameB.StorageLocation != res[1].StorageLocation {
		t.Errorf("game storage location has been changed")
	}

	if gameB.Status != res[1].Status {
		t.Errorf("game status has been changed")
	}

	if gameB.Url != res[1].Url {
		t.Errorf("game url has been changed")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf(err.Error())
	}
}

func Test_Find_Games_When_Database_Is_Empty_Should_Succeed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM games")).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "title", "storage_location", "status", "url"}),
		)

	//Run the test
	repository := repositories.GameRepository(db)

	res, err := repository.FindAll()
	if err != nil {
		t.Errorf(err.Error())
	}

	if res == nil {
		t.Errorf("nil was returned but empty list was expected")
	}

	if len(res) != 0 {
		t.Errorf("FindAll should return an empty list, but it was not empty")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf(err.Error())
	}
}

//************************************ END READ TESTS ************************************
