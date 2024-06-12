package services

import (
	"api/models"
	"api/repositories"
	"api/shared"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"os"
)

var azureBlobContainerName = os.Getenv("AZURE_CONTAINER_NAME")

type IGameService interface {
	FindByID(id uuid.UUID) (*models.Game, error)
	Save(file *multipart.FileHeader, title string, owner string) (*models.Game, error)
	Delete(id uuid.UUID) error
	FindAllByOwner(owner string) ([]models.Game, error)
	ReadOwner(id uuid.UUID) (string, error)
}

type gameService struct {
	repository repositories.IGameRepository
	azClient   *azblob.Client
}

func (g gameService) ReadOwner(id uuid.UUID) (string, error) {
	return g.repository.ReadOwner(id)
}

func (g gameService) FindAllByOwner(owner string) ([]models.Game, error) {
	return g.repository.FindAllByOwner(owner)
}

func (g gameService) FindByID(id uuid.UUID) (*models.Game, error) { return g.repository.FindByID(id) }

func (g gameService) Save(fileHeader *multipart.FileHeader, title string, owner string) (*models.Game, error) {

	game := models.Game{
		ID:              uuid.New(),
		Title:           title,
		StorageLocation: "",
		Status:          shared.Status_New,
		Url:             "",
		Owner:           owner,
	}

	// Creating file on disk because UploadFile() needs *os.File
	dst, err := os.Create(fileHeader.Filename)
	file, err := fileHeader.Open()
	_, err = io.Copy(dst, file)

	if err != nil {
		return nil, err
	}

	_, err = g.azClient.UploadFile(context.Background(), azureBlobContainerName, game.ID.String(), dst, nil)
	if err != nil {
		return nil, err
	}

	// Deleting file from disk
	err = os.Remove(dst.Name())
	if err != nil {
		return nil, err
	}

	storageAccount := os.Getenv("AZURE_STORAGE_ACCOUNT")
	game.StorageLocation = fmt.Sprintf("https://%s.blob.core.windows.net/games/%s", storageAccount, game.ID.String())

	return &game, g.repository.Save(&game)
}

func (g gameService) Delete(id uuid.UUID) error {

	_, err := g.azClient.DeleteBlob(context.Background(), azureBlobContainerName, id.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return g.repository.Delete(id)
}

func GameService(repository repositories.IGameRepository, azClient *azblob.Client) IGameService {
	return &gameService{
		repository: repository,
		azClient:   azClient,
	}
}
