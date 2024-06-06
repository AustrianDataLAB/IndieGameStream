package services

import (
	"api/models"
	"context"
	"errors"
	"github.com/google/uuid"
	streamv1 "indiegamestream.com/indiegamestream/api/stream/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IK8sService interface {
	DeployGame(game *models.Game) error
}

func (g k8sService) DeployGame(game *models.Game) error {

	//Definitions
	resourceName := game.ID.String()
	ctx := context.Background()
	typeNamespacedName := types.NamespacedName{
		Name:      resourceName,
		Namespace: "default",
	}

	//Check if the custom resource is already existing
	err := g.k8sClient.Get(ctx, typeNamespacedName, &streamv1.Game{})
	if err == nil {
		return errors.New("resource is already created")
	} else if !k8serrors.IsNotFound(err) {
		return err
	}

	//Define the custom resource
	resource, err := createAndVerifyGameResource(game)
	if err != nil {
		return err
	}

	return g.k8sClient.Create(ctx, resource)
}

func createAndVerifyGameResource(game *models.Game) (*streamv1.Game, error) {
	if game.ID == uuid.Nil {
		return nil, errors.New("game id is not set")
	}
	if game.Title == "" {
		return nil, errors.New("game title is not set")
	}
	if game.Url == "" {
		return nil, errors.New("game url is not set")
	}

	return &streamv1.Game{
		ObjectMeta: metav1.ObjectMeta{
			Name:      game.ID.String(),
			Namespace: "default",
		},
		Spec: streamv1.GameSpec{
			Name:          game.Title,
			ExecutableURL: game.Url,
		},
	}, nil
}

type k8sService struct {
	k8sClient client.Client
}

func K8sService(k8sClient client.Client) IK8sService {
	return &k8sService{
		k8sClient: k8sClient,
	}
}
