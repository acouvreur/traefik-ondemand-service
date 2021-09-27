package scaler

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

type DockerClassicScaler struct {
	Client client.ContainerAPIClient
}

func NewDockerClassicScaler(client client.ContainerAPIClient) *DockerClassicScaler {
	return &DockerClassicScaler{
		Client: client,
	}
}

func (scaler *DockerClassicScaler) ScaleUp(name string) error {
	log.Infof("Scaling up %s to %d", name, onereplicas)
	ctx := context.Background()
	container, err := scaler.GetContainerByName(name, ctx)

	if err != nil {
		log.Error(err.Error())
		return err
	}

	err = scaler.Client.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})

	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (scaler *DockerClassicScaler) ScaleDown(name string) error {
	log.Infof("Scaling down %s to 0", name)
	ctx := context.Background()
	container, err := scaler.GetContainerByName(name, ctx)

	if err != nil {
		log.Error(err.Error())
		return err
	}

	err = scaler.Client.ContainerStop(ctx, container.ID, nil)

	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (scaler *DockerClassicScaler) IsUp(name string) bool {
	ctx := context.Background()
	container, err := scaler.GetContainerByName(name, ctx)

	if err != nil {
		log.Error(err.Error())
		return false
	}

	spec, err := scaler.Client.ContainerInspect(ctx, container.ID)

	if err != nil {
		log.Error(err.Error())
		return false
	}

	if spec.State.Health != nil {
		return spec.State.Health.Status == "healthy"
	}

	return spec.State.Running
}

func (scaler *DockerClassicScaler) GetContainerByName(name string, ctx context.Context) (*types.Container, error) {
	opts := types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(),
	}
	opts.Filters.Add("name", name)

	containers, err := scaler.Client.ContainerList(ctx, opts)

	if err != nil {
		return nil, err
	}

	if len(containers) == 0 {
		return nil, fmt.Errorf(fmt.Sprintf("container with name %s was not found", name))
	}

	if len(containers) > 1 {
		return nil, fmt.Errorf("multiple containers (%d) with name %s were found: %v", len(containers), name, containers)
	}

	return &containers[0], nil
}
