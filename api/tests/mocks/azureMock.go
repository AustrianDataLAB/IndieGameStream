package mocks

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"io"
	"os"
)

type AzureBlobClientMock struct{}

func (AzureBlobClientMock) CreateContainer(ctx context.Context, containerName string, o *azblob.CreateContainerOptions) (azblob.CreateContainerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (AzureBlobClientMock) DeleteContainer(ctx context.Context, containerName string, o *azblob.DeleteContainerOptions) (azblob.DeleteContainerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (AzureBlobClientMock) DeleteBlob(ctx context.Context, containerName string, blobName string, o *azblob.DeleteBlobOptions) (azblob.DeleteBlobResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (AzureBlobClientMock) NewListBlobsFlatPager(containerName string, o *azblob.ListBlobsFlatOptions) *runtime.Pager[azblob.ListBlobsFlatResponse] {
	//TODO implement me
	panic("implement me")
}

func (AzureBlobClientMock) NewListContainersPager(o *azblob.ListContainersOptions) *runtime.Pager[azblob.ListContainersResponse] {
	//TODO implement me
	panic("implement me")
}

func (AzureBlobClientMock) UploadBuffer(ctx context.Context, containerName string, blobName string, buffer []byte, o *azblob.UploadBufferOptions) (azblob.UploadBufferResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (AzureBlobClientMock) UploadFile(ctx context.Context, containerName string, blobName string, file *os.File, o *azblob.UploadFileOptions) (azblob.UploadFileResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (AzureBlobClientMock) UploadStream(ctx context.Context, containerName string, blobName string, body io.Reader, o *azblob.UploadStreamOptions) (azblob.UploadStreamResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (AzureBlobClientMock) DownloadBuffer(ctx context.Context, containerName string, blobName string, buffer []byte, o *azblob.DownloadBufferOptions) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (AzureBlobClientMock) DownloadFile(ctx context.Context, containerName string, blobName string, file *os.File, o *azblob.DownloadFileOptions) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (AzureBlobClientMock) DownloadStream(ctx context.Context, containerName string, blobName string, o *azblob.DownloadStreamOptions) (azblob.DownloadStreamResponse, error) {
	//TODO implement me
	panic("implement me")
}
