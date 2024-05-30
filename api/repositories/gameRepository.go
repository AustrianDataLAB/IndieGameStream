package repositories

import (
	"api/models"
	"database/sql"
	"github.com/google/uuid"
)

type IGameRepository interface {
	FindAll() ([]models.Game, error)
	FindByID(id uuid.UUID) (*models.Game, error)
	Save(game *models.Game) error
	Delete(id uuid.UUID) error
	FindAllByOwner(owner string) ([]models.Game, error)
	ReadOwner(id uuid.UUID) (string, error)
}

type gameRepository struct {
	db *sql.DB
}

func GameRepository(db *sql.DB) IGameRepository {
	return &gameRepository{
		db: db,
	}
}

// Read the owner of a specific game or empty if the game has not been found
func (g gameRepository) ReadOwner(id uuid.UUID) (string, error) {
	var owner string
	err := g.db.QueryRow("SELECT Owner FROM games WHERE ID = ?", id).Scan(&owner)
	if err != nil {
		return "", err
	}
	return owner, nil
}

// FindAll returns all games from the database or (nil, err) if an error occurred.
func (g gameRepository) FindAll() ([]models.Game, error) {
	query, err := g.db.Query("SELECT * FROM games")
	if err != nil {
		return nil, err
	}
	defer query.Close()
	return readGamesFromRows(query)
}

// FindAll returns all games of a specific owner from the database or (nil, err) if an error occurred.
func (g gameRepository) FindAllByOwner(owner string) ([]models.Game, error) {
	stmt, err := g.db.Prepare("SELECT * FROM games WHERE owner = ?")
	if err != nil {
		return nil, err
	}
	query, err := stmt.Query(owner)
	if err != nil {
		return nil, err
	}
	defer query.Close()
	return readGamesFromRows(query)
}

// FindByID finds a game with a specific id or nil if the game has not been found.
func (g gameRepository) FindByID(id uuid.UUID) (*models.Game, error) {
	var game models.Game
	err := g.db.QueryRow("SELECT * FROM games WHERE ID = ?", id).
		Scan(&game.ID, &game.Title, &game.StorageLocation, &game.Status, &game.Url, &game.Owner)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &game, nil
}

// Save will update the database entry if the game is already in the database.
// If not it will create an uuid and save it in the database.
func (g gameRepository) Save(game *models.Game) error {
	if game.ID != uuid.Nil {
		//Check if uuid is already in database
		existing, err := g.FindByID(game.ID)
		if err != nil {
			return err
		}

		if existing != nil {
			//If yes, update the existing entry
			stmt, err := g.db.Prepare("UPDATE games SET Title=?, StorageLocation=?, Status=?, Url=? WHERE ID = ?")
			if err != nil {
				return err
			}

			return checkResult(stmt.Exec(game.Title, game.StorageLocation, game.Status, game.Url, game.ID))
		}
	} else {
		game.ID = uuid.New()
	}

	//If not create a new one
	stmt, err := g.db.Prepare("INSERT INTO games (ID, Title, StorageLocation, Status, Url, Owner) VALUES (?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	return checkResult(stmt.Exec(game.ID, game.Title, game.StorageLocation, game.Status, game.Url, game.Owner))
}

// Delete removes the entry with a specific id from the games database.
// Or returns sql.ErrNoRows if the game is not existing.
func (g gameRepository) Delete(id uuid.UUID) error {
	stmt, err := g.db.Prepare("DELETE FROM games WHERE ID = ?")
	if err != nil {
		return err
	}

	return checkResult(stmt.Exec(id))
}

func checkResult(res sql.Result, err error) error {
	if err != nil {
		return err
	}

	return checkAffectedRows(res)
}

func checkAffectedRows(res sql.Result) error {
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func readGamesFromRows(query *sql.Rows) ([]models.Game, error) {
	var games = []models.Game{}
	for query.Next() {
		var game models.Game
		err := query.Scan(&game.ID, &game.Title, &game.StorageLocation, &game.Status, &game.Url, &game.Owner)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	err := query.Err()
	if err != nil {
		return nil, err
	}

	return games, nil
}
