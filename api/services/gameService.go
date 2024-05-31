package services

import (
	"api/models"
	"api/repositories"
	"api/shared"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type IGameService interface {
	FindAll() ([]models.Game, error)
	FindByID(id uuid.UUID) (*models.Game, error)
	Save(file *multipart.FileHeader, title string) (*models.Game, error)
	Delete(id uuid.UUID) error
}

type gameService struct {
	repository repositories.IGameRepository
}

func (g gameService) FindAll() ([]models.Game, error) {
	return g.repository.FindAll()
}

func (g gameService) FindByID(id uuid.UUID) (*models.Game, error) { return g.repository.FindByID(id) }

func (g gameService) Save(file *multipart.FileHeader, title string) (*models.Game, error) {
	//TODO save file

	game := models.Game{
		ID:              uuid.New(),
		Title:           title,
		StorageLocation: "",
		Status:          shared.Status_New,
		Url:             "",
	}

	// TODO replace stamango
	url := "https://stamango.blob.core.windows.net/"

	//clientID := os.Getenv("AZURE_CLIENT_ID")
	//opts := &azidentity.ManagedIdentityCredentialOptions{ID: azidentity.ClientID(clientID)}
	//credential, err := azidentity.NewManagedIdentityCredential(opts)
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Println(fmt.Sprintf("------------------------- err1: %s", err))
		return nil, err
	}

	client, err := azblob.NewClient(url, credential, nil)
	if err != nil {
		return nil, err
	}

	containerName := "indiegamestream"
	_, err = client.CreateContainer(context.Background(), containerName, nil)
	if err != nil {
		return nil, err
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
