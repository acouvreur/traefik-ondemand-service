package scaler

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

type DockerSwarmScaler struct{}

func (DockerSwarmScaler) ScaleUp(client *client.Client, name string, replicas *uint64) {
	log.Infof("Scaling up %s to %d", name, *replicas)
	ctx := context.Background()
	service, err := GetServiceByName(client, name, ctx)

	if err != nil {
		println(err)
		return
	}

	service.Spec.Mode.Replicated = &swarm.ReplicatedService{
		Replicas: replicas,
	}
	response, err := client.ServiceUpdate(ctx, service.ID, service.Meta.Version, service.Spec, types.ServiceUpdateOptions{})

	if err != nil {
		println(err)
		return
	}

	if len(response.Warnings) > 0 {
		fmt.Printf("Warnings received scaling up service %s: %v", name, response.Warnings)
	}
}

func (DockerSwarmScaler) ScaleDown(client *client.Client, name string) {
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

func (DockerSwarmScaler) IsUp(client *client.Client, name string) bool {
	ctx := context.Background()
	service, err := GetServiceByName(client, name, ctx)

	if err != nil {
		println(err)
		return false
	}

	return *service.Spec.Mode.Replicated.Replicas > 0
}

func GetServiceByName(client *client.Client, name string, ctx context.Context) (*swarm.Service, error) {
	opts := types.ServiceListOptions{
		Filters: filters.NewArgs(),
	}
	opts.Filters.Add("name", name)

	services, err := client.ServiceList(ctx, opts)

	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, fmt.Errorf(fmt.Sprintf("service with name %s was not found", name))
	}

	if len(services) > 1 {
		return nil, fmt.Errorf("multiple services (%d) with name %s were found: %v", len(services), name, services)
	}

	return &services[0], nil
}
