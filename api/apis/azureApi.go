package apis

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"io"
	"mime/multipart"
	"os"
)

var storageAccount = os.Getenv("AZURE_STORAGE_ACCOUNT")

type IAzureApi interface {
	UploadGame(blobContainerName string, gameID string, fileHeader *multipart.FileHeader) (string, error)
	DeleteGame(blobContainerName string, gameID string) error
}

func (g azureApi) UploadGame(blobContainerName string, gameID string, fileHeader *multipart.FileHeader) (string, error) {
	ctx := context.Background()

	// Creating file on disk because UploadFile() needs *os.File
	dst, err := os.Create(fileHeader.Filename)
	file, err := fileHeader.Open()
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	_, err = g.azure.UploadFile(ctx, blobContainerName, gameID, dst, nil)
	if err != nil {
		return "", err
	}

	// Deleting file from disk
	err = os.Remove(dst.Name())
	if err != nil {
		return "", err
	}

	storageLocation := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", storageAccount, blobContainerName, gameID)

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
