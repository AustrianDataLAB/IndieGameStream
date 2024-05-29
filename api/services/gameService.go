package services

import (
	"api/models"
	"api/repositories"
	"api/shared"
	"github.com/google/uuid"
	"mime/multipart"
)

type IGameService interface {
	FindAll() ([]models.Game, error)
	FindByID(id uuid.UUID) (*models.Game, error)
	Save(file *multipart.FileHeader, title string, owner string) (*models.Game, error)
	Delete(id uuid.UUID) error
}

type gameService struct {
	repository repositories.IGameRepository
}

func (g gameService) FindAll() ([]models.Game, error) {
	return g.repository.FindAll()
}

func (g gameService) FindByID(id uuid.UUID) (*models.Game, error) { return g.repository.FindByID(id) }

func (g gameService) Save(file *multipart.FileHeader, title string, owner string) (*models.Game, error) {
	//TODO save file

	game := models.Game{
		ID:              uuid.New(),
		Title:           title,
		StorageLocation: "",
		Status:          shared.Status_New,
		Url:             "",
		Owner:           owner,
	}
	return &game, g.repository.Save(&game)
}

func (g gameService) Delete(id uuid.UUID) error {
	return g.repository.Delete(id)
}

func GameService(repository repositories.IGameRepository) IGameService {
	return &gameService{
		repository: repository,
	}
}
