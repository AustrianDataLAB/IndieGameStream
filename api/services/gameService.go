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
	}

	//If the game url is empty, try to get it from kubernetes
	if game.Url == "" {
		g.updateGameUrl(game)
	}

	return g.repository.FindByID(id)

}

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

	//**********************************************
	//**** Placeholder for upload game to azure ****
	//**********************************************

	//Deploy the game on kubernetes
	err := g.k8s.DeployGame(&game)
	if err != nil {
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
	return g.repository.Delete(id)
}

func (g gameService) updateGameUrl(game *models.Game) {
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

func GameService(repository repositories.IGameRepository, k8s apis.IK8sApi) IGameService {
	return &gameService{
		repository: repository,
		k8s:        k8s,
	}
}
