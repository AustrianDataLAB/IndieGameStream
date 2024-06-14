package apis

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"io"
	"mime/multipart"
	"os"
)

type IAzureApi interface {
	UploadGame(blobContainerName string, gameID string, fileHeader *multipart.FileHeader) (string, error)
	DeleteGame(blobContainerName string, gameID string) error
}

func (g azureApi) UploadGame(blobContainerName string, gameID string, fileHeader *multipart.FileHeader) (string, error) {
	ctx := context.Background()

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}

	gameBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	_, err = g.azure.UploadBuffer(ctx, blobContainerName, gameID, gameBytes, nil)
	if err != nil {
		return "", err
	}

	storageLocation := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", os.Getenv("AZURE_STORAGE_ACCOUNT"), blobContainerName, gameID)

	return storageLocation, nil
}

func (g azureApi) DeleteGame(blobContainerName string, gameID string) error {
	ctx := context.Background()

	_, err := g.azure.DeleteBlob(ctx, blobContainerName, gameID, nil)
	if err != nil {
		return err
	}

	return nil
}

type azureApi struct {
	azure *azblob.Client
}

func AzureService(azure *azblob.Client) IAzureApi {
	return &azureApi{
		azure: azure,
	}
}
