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
		Title:           "MockTitle",
		StorageLocation: "MockStorageLocation",
		Status:          "MockStatus",
		Url:             "MockUrl",
		Owner:           "MockOwner",
		FileName:        "TestFile.nes",
	}
	mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO games"))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO games")).
		WithArgs(sqlmock.AnyArg(), game.Title, game.StorageLocation, game.Status, game.Url, game.Owner, game.FileName).
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
		Title:           "MockTitle",
		StorageLocation: "MockStorageLocation",
		Status:          "MockStatus",
		Url:             "MockUrl",
		Owner:           "MockOwner",
		FileName:        "TestFile.nes",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM games WHERE ID = ?")).
		WithArgs(id).WillReturnError(sql.ErrNoRows)

	mock.ExpectPrepare("INSERT INTO games")

	mock.ExpectExec("INSERT INTO games").
		WithArgs(game.ID, game.Title, game.StorageLocation, game.Status, game.Url, game.Owner, game.FileName).
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
		Title:           "MockTitle",
		StorageLocation: "MockStorageLocation",
		Status:          "MockStatus",
		Url:             "MockUrl",
		Owner:           "MockOwner",
		FileName:        "TestFile.nes",
	}
	
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM games WHERE ID = ?")).
		WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"Id"}).AddRow(id))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM games WHERE ID = ?")).
		WithArgs(id).WillReturnRows(
		sqlmock.NewRows([]string{"Id", "Title", "StorageLocation", "Status", "Url", "Owner", "FileName"}).
			AddRow(id, "", "", "", "", "", ""),
	)

	mock.ExpectPrepare(regexp.
		QuoteMeta("UPDATE games SET Title=?, StorageLocation=?, Status=?, Url=?, FileName=? WHERE ID = ?"))

	mock.ExpectExec(regexp.
		QuoteMeta("UPDATE games SET Title=?, StorageLocation=?, Status=?, Url=?, FileName=? WHERE ID = ?")).
		WithArgs(game.Title, game.StorageLocation, game.Status, game.Url, game.ID, game.FileName).
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
		Title:           "MockTitle",
		StorageLocation: "MockStorageLocation",
		Status:          "MockStatus",
		Url:             "MockUrl",
		Owner:           "MockOwner",
		FileName:        "TestFile.nes",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM games WHERE ID = ?")).
		WithArgs(id).WillReturnRows(
		sqlmock.NewRows([]string{"Id", "Title", "StorageLocation", "Status", "Url", "Owner", "FileName"}).
			AddRow(id, game.Title, game.StorageLocation, game.Status, game.Url, game.Owner, game.FileName),
	)

	//Run the test
	repository := repositories.GameRepository(db)

	res, err := repository.FindByID(game.ID)
	if err != nil {
		t.Errorf(err.Error())
	}

	compareGames(t, &game, res)

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

func Test_Find_Two_Games_Of_Same_Owner_Should_Succeed(t *testing.T) {
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
		Owner:           "MockOwner",
		FileName:        "TestFile.nes",
	}

	gameB := models.Game{
		ID:              uuid.New(),
		Title:           "BMockA",
		StorageLocation: "BMockB",
		Status:          "BMockC",
		Url:             "BMockD",
		Owner:           "MockOwner",
		FileName:        "TestFile2.nes",
	}

	mock.ExpectPrepare(regexp.QuoteMeta("SELECT * FROM games WHERE owner = ?"))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM games WHERE owner = ?")).
		WithArgs("MockOwner").
		WillReturnRows(
			sqlmock.NewRows([]string{"ID", "Title", "StorageLocation", "Status", "Url", "Owner", "FileName"}).
				AddRow(gameA.ID, gameA.Title, gameA.StorageLocation, gameA.Status, gameA.Url, gameA.Owner, gameA.FileName).
				AddRow(gameB.ID, gameB.Title, gameB.StorageLocation, gameB.Status, gameB.Url, gameB.Owner, gameB.FileName),
		)

	//Run the test
	repository := repositories.GameRepository(db)

	res, err := repository.FindAllByOwner("MockOwner")
	if err != nil {
		t.Errorf(err.Error())
	}

	if res == nil {
		t.Errorf("game was not returned")
	}

	if len(res) != 2 {
		t.Errorf("FindAll should return two games")
	}

	compareGames(t, &gameA, &res[0])
	compareGames(t, &gameB, &res[1])

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

	mock.ExpectPrepare(regexp.QuoteMeta("SELECT * FROM games WHERE owner = ?"))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM games WHERE owner = ?")).
		WithArgs("MockOwner").
		WillReturnRows(
			sqlmock.NewRows([]string{"Id", "Title", "StorageLocation", "Status", "Url", "Owner", "FileName"}),
		)

	//Run the test
	repository := repositories.GameRepository(db)

	res, err := repository.FindAllByOwner("MockOwner")
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

func Test_Read_Owner_Should_Succeed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer db.Close()

	id := uuid.New()
	game := models.Game{
		ID:              id,
		Title:           "MockTitle",
		StorageLocation: "MockStorageLocation",
		Status:          "MockStatus",
		Url:             "MockUrl",
		Owner:           "MockOwner",
		FileName:        "TestFile.nes",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT Owner FROM games WHERE ID = ?")).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"Id", "Title", "StorageLocation", "Status", "Url", "Owner", "FileName"}))

	//Run the test
	repository := repositories.GameRepository(db)

	res, err := repository.ReadOwner(game.ID)
	if err == nil {
		t.Errorf("error was not returned")
	}

	if res != "" {
		t.Errorf("Owner should be empty, but got %s", res)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf(err.Error())
	}
}

func Test_Read_Owner_Should_Throw_When_Database_Is_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer db.Close()

	id := uuid.New()
	game := models.Game{
		ID:              id,
		Title:           "MockTitle",
		StorageLocation: "MockStorageLocation",
		Status:          "MockStatus",
		Url:             "MockUrl",
		Owner:           "MockOwner",
		FileName:        "TestFile.nes",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT Owner FROM games WHERE ID = ?")).
		WithArgs(id).WillReturnRows(
		sqlmock.NewRows([]string{"owner"}).
			AddRow(game.Owner),
	)

	//Run the test
	repository := repositories.GameRepository(db)

	res, err := repository.ReadOwner(game.ID)
	if err != nil {
		t.Errorf(err.Error())
	}

	if res != game.Owner {
		t.Errorf("ReadOwner should return %v, but got %v", game.Owner, res)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf(err.Error())
	}
}

//************************************ END READ TESTS ************************************

func compareGames(t *testing.T, game *models.Game, res *models.Game) {
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

	if game.Owner != res.Owner {
		t.Errorf("game owner has been changed")
	}

	if game.FileName != res.FileName {
		t.Errorf("game file name has been changed")
	}

}
