package azureClient

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"io"
	"os"
)

// azblob.client has no interface
// but we need one for mocking in the tests.
// So we have to define our own interface

type IAzBlobClient interface {
	CreateContainer(ctx context.Context, containerName string, o *azblob.CreateContainerOptions) (azblob.CreateContainerResponse, error)
	DeleteContainer(ctx context.Context, containerName string, o *azblob.DeleteContainerOptions) (azblob.DeleteContainerResponse, error)
	DeleteBlob(ctx context.Context, containerName string, blobName string, o *azblob.DeleteBlobOptions) (azblob.DeleteBlobResponse, error)
	NewListBlobsFlatPager(containerName string, o *azblob.ListBlobsFlatOptions) *runtime.Pager[azblob.ListBlobsFlatResponse]
	NewListContainersPager(o *azblob.ListContainersOptions) *runtime.Pager[azblob.ListContainersResponse]
	UploadBuffer(ctx context.Context, containerName string, blobName string, buffer []byte, o *azblob.UploadBufferOptions) (azblob.UploadBufferResponse, error)
	UploadFile(ctx context.Context, containerName string, blobName string, file *os.File, o *azblob.UploadFileOptions) (azblob.UploadFileResponse, error)
	UploadStream(ctx context.Context, containerName string, blobName string, body io.Reader, o *azblob.UploadStreamOptions) (azblob.UploadStreamResponse, error)
	DownloadBuffer(ctx context.Context, containerName string, blobName string, buffer []byte, o *azblob.DownloadBufferOptions) (int64, error)
	DownloadFile(ctx context.Context, containerName string, blobName string, file *os.File, o *azblob.DownloadFileOptions) (int64, error)
	DownloadStream(ctx context.Context, containerName string, blobName string, o *azblob.DownloadStreamOptions) (azblob.DownloadStreamResponse, error)
}

type azBlobClient struct {
	client *azblob.Client
}

func AzBlobClient(azure *azblob.Client) IAzBlobClient {
	return &azBlobClient{
		client: azure,
	}
}

// URL returns the URL endpoint used by the BlobClient object.
func (c *azBlobClient) URL() string {
	return c.URL()
}

// ServiceClient returns the embedded service client for this client.
func (c *azBlobClient) ServiceClient() *service.Client {
	return c.client.ServiceClient()
}

// CreateContainer is a lifecycle method to creates a new container under the specified account.
// If the container with the same name already exists, a ResourceExistsError will be raised.
// This method returns a client with which to interact with the newly created container.
func (c *azBlobClient) CreateContainer(ctx context.Context, containerName string, o *azblob.CreateContainerOptions) (azblob.CreateContainerResponse, error) {
	return c.client.CreateContainer(ctx, containerName, o)
}

// DeleteContainer is a lifecycle method that marks the specified container for deletion.
// The container and any blobs contained within it are later deleted during garbage collection.
// If the container is not found, a ResourceNotFoundError will be raised.
func (c *azBlobClient) DeleteContainer(ctx context.Context, containerName string, o *azblob.DeleteContainerOptions) (azblob.DeleteContainerResponse, error) {
	return c.client.DeleteContainer(ctx, containerName, o)
}

// DeleteBlob marks the specified blob or snapshot for deletion. The blob is later deleted during garbage collection.
// Note that deleting a blob also deletes all its snapshots.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/delete-blob.
func (c *azBlobClient) DeleteBlob(ctx context.Context, containerName string, blobName string, o *azblob.DeleteBlobOptions) (azblob.DeleteBlobResponse, error) {
	return c.client.DeleteBlob(ctx, containerName, blobName, o)
}

// NewListBlobsFlatPager returns a pager for blobs starting from the specified Marker. Use an empty
// Marker to start enumeration from the beginning. Blob names are returned in lexicographic order.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/list-blobs.
func (c *azBlobClient) NewListBlobsFlatPager(containerName string, o *azblob.ListBlobsFlatOptions) *runtime.Pager[azblob.ListBlobsFlatResponse] {
	return c.client.NewListBlobsFlatPager(containerName, o)
}

// NewListContainersPager operation returns a pager of the containers under the specified account.
// Use an empty Marker to start enumeration from the beginning. Container names are returned in lexicographic order.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/list-containers2.
func (c *azBlobClient) NewListContainersPager(o *azblob.ListContainersOptions) *runtime.Pager[azblob.ListContainersResponse] {
	return c.client.NewListContainersPager(o)
}

// UploadBuffer uploads a buffer in blocks to a block blob.
func (c *azBlobClient) UploadBuffer(ctx context.Context, containerName string, blobName string, buffer []byte, o *azblob.UploadBufferOptions) (azblob.UploadBufferResponse, error) {
	return c.client.UploadBuffer(ctx, containerName, blobName, buffer, o)
}

// UploadFile uploads a file in blocks to a block blob.
func (c *azBlobClient) UploadFile(ctx context.Context, containerName string, blobName string, file *os.File, o *azblob.UploadFileOptions) (azblob.UploadFileResponse, error) {
	return c.client.UploadFile(ctx, containerName, blobName, file, o)
}

// UploadStream copies the file held in io.Reader to the Blob at blockBlobClient.
// A Context deadline or cancellation will cause this to error.
func (c *azBlobClient) UploadStream(ctx context.Context, containerName string, blobName string, body io.Reader, o *azblob.UploadStreamOptions) (azblob.UploadStreamResponse, error) {
	return c.client.UploadStream(ctx, containerName, blobName, body, o)
}

// DownloadBuffer downloads an Azure blob to a buffer with parallel.
func (c *azBlobClient) DownloadBuffer(ctx context.Context, containerName string, blobName string, buffer []byte, o *azblob.DownloadBufferOptions) (int64, error) {
	return c.client.DownloadBuffer(ctx, containerName, blobName, buffer, o)
}

// DownloadFile downloads an Azure blob to a local file.
// The file would be truncated if the size doesn't match.
func (c *azBlobClient) DownloadFile(ctx context.Context, containerName string, blobName string, file *os.File, o *azblob.DownloadFileOptions) (int64, error) {
	return c.client.DownloadFile(ctx, containerName, blobName, file, o)
}

// DownloadStream reads a range of bytes from a blob. The response also includes the blob's properties and metadata.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-blob.
func (c *azBlobClient) DownloadStream(ctx context.Context, containerName string, blobName string, o *azblob.DownloadStreamOptions) (azblob.DownloadStreamResponse, error) {
	return c.client.DownloadStream(ctx, containerName, blobName, o)
}
