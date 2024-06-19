package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type k8sMock struct {
	mock *mock.Mock
}

func K8sMock(mock *mock.Mock) *k8sMock {
	return &k8sMock{mock: mock}
}

func (s k8sMock) Mock() *mock.Mock {
	return s.mock
}

func (s k8sMock) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	s.mock.Called(ctx, key, obj)
	return nil
}

func (s k8sMock) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	s.mock.Called(ctx, list)
	return nil
}

func (s k8sMock) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	s.mock.Called(ctx, obj)
	return nil
}

func (s k8sMock) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	s.mock.Called(ctx, obj)
	return nil
}

func (s k8sMock) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	s.mock.Called(ctx, obj)
	return nil
}

func (s k8sMock) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	s.mock.Called(ctx, obj, patch)
	return nil
}

func (s k8sMock) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	s.mock.Called(ctx, obj)
	return nil
}

func (s k8sMock) Status() client.SubResourceWriter {
	s.mock.Called()
	return nil
}

func (s k8sMock) SubResource(subResource string) client.SubResourceClient {
	s.mock.Called(subResource)
	return nil
}

func (s k8sMock) Scheme() *runtime.Scheme {
	s.mock.Called()
	return nil
}

func (s k8sMock) RESTMapper() meta.RESTMapper {
	s.mock.Called()
	return nil
}

func (s k8sMock) GroupVersionKindFor(obj runtime.Object) (schema.GroupVersionKind, error) {
	s.mock.Called(obj)
	return schema.GroupVersionKind{}, nil
}

func (s k8sMock) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	s.mock.Called(obj)
	return true, nil
}
