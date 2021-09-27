package mocks

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/mock"
)

type ContainerAPIClientMock struct {
	client.ContainerAPIClient
	mock.Mock
}

func NewContainerAPIClientMock() *ContainerAPIClientMock {
	return &ContainerAPIClientMock{}
}

func (client *ContainerAPIClientMock) ContainerStart(ctx context.Context, container string, options types.ContainerStartOptions) error {
	args := client.Mock.Called(ctx, container, options)
	return args.Error(0)
}

func (client *ContainerAPIClientMock) ContainerStop(ctx context.Context, container string, timeout *time.Duration) error {
	args := client.Mock.Called(ctx, container, timeout)
	return args.Error(0)
}

func (client *ContainerAPIClientMock) ContainerInspect(ctx context.Context, container string) (types.ContainerJSON, error) {
	args := client.Mock.Called(ctx, container)
	return args.Get(0).(types.ContainerJSON), args.Error(1)
}

func (client *ContainerAPIClientMock) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	args := client.Mock.Called(ctx, options)
	return args.Get(0).([]types.Container), args.Error(1)
}

type ServiceAPIClientMock struct {
	client.ServiceAPIClient
	mock.Mock
}

func NewServiceAPIClientMock() *ServiceAPIClientMock {
	return &ServiceAPIClientMock{}
}

func (client *ServiceAPIClientMock) ServiceUpdate(ctx context.Context, serviceID string, version swarm.Version, service swarm.ServiceSpec, options types.ServiceUpdateOptions) (types.ServiceUpdateResponse, error) {
	args := client.Mock.Called(ctx, serviceID, version, service, options)
	return args.Get(0).(types.ServiceUpdateResponse), args.Error(1)
}

func (client *ServiceAPIClientMock) ServiceList(ctx context.Context, options types.ServiceListOptions) ([]swarm.Service, error) {
	args := client.Mock.Called(ctx, options)
	return args.Get(0).([]swarm.Service), args.Error(1)
}
