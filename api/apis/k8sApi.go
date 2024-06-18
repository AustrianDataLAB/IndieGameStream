package apis

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

type IK8sApi interface {
	DeployGame(game *models.Game) error
	ReadGameUrl(gameId uuid.UUID) (string, error)
	DeleteGame(game *models.Game) error
}

func (g k8sApi) DeleteGame(game *models.Game) error {
	resource, err := createAndVerifyGameResource(game)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return g.k8sClient.Delete(ctx, resource)
}

func (g k8sApi) DeployGame(game *models.Game) error {

	//Definitions
	ctx := context.Background()
	key := typeNamespacedName(game.ID.String())

	//Check if the custom resource is already existing
	err := g.k8sClient.Get(ctx, key, &streamv1.Game{})
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

func (g k8sApi) ReadGameUrl(gameId uuid.UUID) (string, error) {
	key := typeNamespacedName(gameId.String())
	resource := streamv1.Game{}

	err := g.k8sClient.Get(context.Background(), key, &resource)
	if err != nil {
		return "", err
	} else {
		return resource.Status.URL, nil
	}
}

func typeNamespacedName(resourceName string) types.NamespacedName {
	return types.NamespacedName{
		Name:      resourceName,
		Namespace: "default",
	}
}

func createAndVerifyGameResource(game *models.Game) (*streamv1.Game, error) {
	if game.ID == uuid.Nil {
		return nil, errors.New("game id is not set")
	}
	if game.Title == "" {
		return nil, errors.New("game title is not set")
	}
	if game.StorageLocation == "" {
		return nil, errors.New("game StorageLocation is not set")
	}

	return &streamv1.Game{
		ObjectMeta: metav1.ObjectMeta{
			Name:      game.ID.String(),
			Namespace: "default",
		},
		Spec: streamv1.GameSpec{
			Name:     game.Title,
			FileName: game.FileName,
		},
	}, nil
}

type k8sApi struct {
	k8sClient client.Client
}

func K8sService(k8sClient client.Client) IK8sApi {
	return &k8sApi{
		k8sClient: k8sClient,
	}
}
