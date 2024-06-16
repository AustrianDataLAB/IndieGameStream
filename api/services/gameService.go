package services

import (
	"api/apis"
	"api/models"
	"api/repositories"
	"api/shared"
	"fmt"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"os"
	"strings"
)

type IGameService interface {
	FindByID(id uuid.UUID) (*models.Game, error)
	Save(file *multipart.FileHeader, title string, owner string) (*models.Game, error)
	Delete(id uuid.UUID) error
	FindAllByOwner(owner string) ([]models.Game, error)
	ReadOwner(id uuid.UUID) (string, error)
}

type gameService struct {
	repository repositories.IGameRepository
	azure      apis.IAzureApi
	k8s        apis.IK8sApi
}

func (g gameService) ReadOwner(id uuid.UUID) (string, error) {
	return g.repository.ReadOwner(id)
}

func (g gameService) FindAllByOwner(owner string) ([]models.Game, error) {
	return g.repository.FindAllByOwner(owner)
}

func (g gameService) FindByID(id uuid.UUID) (*models.Game, error) {
	game, err := g.repository.FindByID(id)
	if err != nil {
		return nil, err
	} else {
		g.updateGameUrl(game)
		return game, nil
	}
}

func (g gameService) Save(fileHeader *multipart.FileHeader, title string, owner string) (*models.Game, error) {

	game := models.Game{
		ID:              uuid.New(),
		Title:           title,
		StorageLocation: "",
		Status:          shared.Status_New,
		Url:             "",
		Owner:           owner,
	}

	//Upload game to azure blob storage container
	storageLocation, err := g.azure.UploadGame(os.Getenv("AZURE_CONTAINER_NAME"), game.ID.String(), fileHeader)
	if err != nil {
		return nil, err
	}

	game.StorageLocation = storageLocation

	//Deploy the game on kubernetes
	err = g.k8s.DeployGame(&game)
	if err != nil {
		//Delete the game when deploying on kubernetes failed
		errDel := g.azure.DeleteGame(os.Getenv("AZURE_CONTAINER_NAME"), game.ID.String())
		if errDel != nil {
			log.Println(fmt.Sprintf("Delete game for %s in azure failed", title))
		}

		return nil, err
	}

	//Try to read the game url
	game.Url, err = g.k8s.ReadGameUrl(game.ID)
	if err != nil {
		log.Println(fmt.Sprintf("Error reading game url: %s", err))
		//We can ignore this error because we try it again in FindByID
		//Maybe the deployment is not ready yet
	}

	return &game, g.repository.Save(&game)
}

func (g gameService) Delete(id uuid.UUID) error {
	//Get the game, we need the details to delete it from k8s
	game, err := g.repository.FindByID(id)
	if err != nil {
		return err
	}

	//Delete from azure storage
	err = g.azure.DeleteGame(os.Getenv("AZURE_CONTAINER_NAME"), id.String())
	if err != nil {
		if isNotFound(err) {
			log.Println(fmt.Sprintf("Game %s is already deleted from azure storage", id.String()))
		} else {
			return err
		}
	}

	//Delete from k8s/aks, if the game has an url
	if game.Url != "" {
		err = g.k8s.DeleteGame(game)
		if err != nil {
			if isNotFound(err) {
				log.Println(fmt.Sprintf("Game %s is already deleted from aks", id.String()))
			} else {
				return err
			}
		}
	}

	//Delete from db and return
	return g.repository.Delete(id)
}

func (g gameService) updateGameUrl(game *models.Game) {
	if game == nil {
		return
	}

	url, err := g.k8s.ReadGameUrl(game.ID)
	if err != nil {
		log.Println(fmt.Sprintf("Error reading game url: %s", err))
		//We can ignore this error because we try it again next time
		//Maybe the deployment is not ready yet
	} else {
		game.Url = url
		//Save the changes in the database
		err := g.repository.Save(game)
		if err != nil {
			log.Println(fmt.Sprintf("Error updating game: %s", err))
		}
	}
}

func GameService(repository repositories.IGameRepository, k8s apis.IK8sApi, azure apis.IAzureApi) IGameService {
	return &gameService{
		repository: repository,
		k8s:        k8s,
		azure:      azure,
	}
}

func isNotFound(err error) bool {
	e := strings.ToLower(err.Error())
	return strings.Contains(e, "not found") || strings.Contains(e, "notfound")
}
