package repositories

import (
	"api/models"
	"api/shared"
	"github.com/google/uuid"
)

type IGameRepository interface {
	FindAll() ([]models.Game, error)
	FindByID(id uuid.UUID) (models.Game, error)
	Save(game models.Game) (models.Game, error)
	Delete(id uuid.UUID) error
}

type gameRepository struct {
	//TODO add database
}

func GameRepository() IGameRepository {
	return &gameRepository{
		//TODO add database
	}
}

func (g gameRepository) FindAll() ([]models.Game, error) {
	//TODO implement me
	return []models.Game{
		{
			ID:              uuid.New(),
			Title:           "Mock1",
			StorageLocation: "Mock1",
			Status:          shared.Status_Installed,
			Url:             "https://localhost:4200",
		},
		{
			ID:              uuid.New(),
			Title:           "Mock2",
			StorageLocation: "Mock2",
			Status:          shared.Status_Installed,
			Url:             "https://localhost:4200",
		},
	}, nil
}

func (g gameRepository) FindByID(id uuid.UUID) (models.Game, error) {
	//TODO implement me
	return models.Game{
		ID:              id,
		Title:           "Mock",
		StorageLocation: "Mock",
		Status:          shared.Status_Installed,
		Url:             "https://localhost:4200",
	}, nil
}

func (g gameRepository) Save(game models.Game) (models.Game, error) {
	//TODO implement me
	game.ID = uuid.New()
	return game, nil
}

func (g gameRepository) Delete(id uuid.UUID) error {
	//TODO implement me
	return nil
}
