package services

import (
	"api/models"
	"api/repositories"
	"api/shared"
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"os"
)

type IGameService interface {
	FindAll() ([]models.Game, error)
	FindByID(id uuid.UUID) (*models.Game, error)
	Save(file *multipart.FileHeader, title string) (*models.Game, error)
	Delete(id uuid.UUID) error
}

type gameService struct {
	repository repositories.IGameRepository
	azClient *azblob.Client
}

func (g gameService) FindAll() ([]models.Game, error) {
	return g.repository.FindAll()
}

func (g gameService) FindByID(id uuid.UUID) (*models.Game, error) { return g.repository.FindByID(id) }

func (g gameService) Save(fileHeader *multipart.FileHeader, title string) (*models.Game, error) {

	game := models.Game{
		ID:              uuid.New(),
		Title:           title,
		StorageLocation: "",
		Status:          shared.Status_New,
		Url:             "",
	}

	containerName := os.Getenv("AZURE_CONTAINER_NAME")

	// Creating file on disk because UploadFile() needs *os.File
	dst, err := os.Create(fileHeader.Filename)
	file, err := fileHeader.Open()
	_, err = io.Copy(dst, file)

	if err != nil {
		log.Fatal(err)
	}

	// TODO RESPONSE 403: This request is not authorized to perform this operation using this permission. ERROR CODE: AuthorizationPermissionMismatch
	_, err = g.azClient.UploadFile(context.Background(), containerName, game.ID.String(), dst, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Deleting file from disk
	err = os.Remove(dst.Name())
	if err != nil {
		log.Fatal(err)
	}

	storageAccount := os.Getenv("AZURE_STORAGE_ACCOUNT")
	game.StorageLocation = fmt.Sprintf("https://%s.blob.core.windows.net/games/%s", storageAccount, game.ID.String())

	return &game, g.repository.Save(&game)
}

func (g gameService) Delete(id uuid.UUID) error {
	return g.repository.Delete(id)
}

func GameService(repository repositories.IGameRepository, azClient *azblob.Client) IGameService {
	return &gameService{
		repository: repository,
		azClient: azClient,
	}
}
