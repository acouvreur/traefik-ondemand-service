package scaler

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

type DockerClassicScaler struct{}

func (DockerClassicScaler) ScaleUp(client *client.Client, name string, replicas *uint64) {
	log.Infof("Scaling up %s to %d", name, *replicas)
	ctx := context.Background()
	container, err := GetContainerByName(client, name, ctx)

	if err != nil {
		println(err)
		return
	}

	err = client.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})

	if err != nil {
		println(err)
		return
	}
}

func (DockerClassicScaler) ScaleDown(client *client.Client, name string) {
	log.Infof("Scaling down %s to 0", name)
	ctx := context.Background()
	container, err := GetContainerByName(client, name, ctx)

	if err != nil {
		println(err)
		return
	}

	err = client.ContainerStop(ctx, container.ID, nil)

	if err != nil {
		println(err)
		return
	}
}

func (DockerClassicScaler) IsUp(client *client.Client, name string) bool {
	ctx := context.Background()
	container, err := GetContainerByName(client, name, ctx)

	if err != nil {
		println(err)
		return false
	}

	spec, err := client.ContainerInspect(ctx, container.ID)

	if err != nil {
		println(err)
		return false
	}

	return spec.State.Running
}

func GetContainerByName(client *client.Client, name string, ctx context.Context) (*types.Container, error) {
	opts := types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(),
	}
	opts.Filters.Add("name", name)

	containers, err := client.ContainerList(ctx, opts)

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
