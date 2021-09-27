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

type DockerSwarmScaler struct {
	Client client.ServiceAPIClient
}

func NewDockerSwarmScaler() *DockerSwarmScaler {
	return &DockerSwarmScaler{}
}

func (scaler *DockerSwarmScaler) ScaleUp(name string) error {
	log.Infof("scaling up %s to %d", name, onereplicas)
	ctx := context.Background()
	service, err := scaler.GetServiceByName(name, ctx)

	if err != nil {
		log.Error(err.Error())
		return err
	}

	service.Spec.Mode.Replicated = &swarm.ReplicatedService{
		Replicas: &onereplicas,
	}
	response, err := scaler.Client.ServiceUpdate(ctx, service.ID, service.Meta.Version, service.Spec, types.ServiceUpdateOptions{})

	if err != nil {
		log.Error(err.Error())
		return err
	}

	if len(response.Warnings) > 0 {
		log.Warnf("received scaling up service %s: %v", name, response.Warnings)
	}

	return nil
}

func (scaler *DockerSwarmScaler) ScaleDown(name string) error {
	log.Infof("scaling down %s to 0", name)
	ctx := context.Background()
	service, err := scaler.GetServiceByName(name, ctx)

	if err != nil {
		log.Error(err.Error())
		return err
	}

	replicas := uint64(0)

	service.Spec.Mode.Replicated = &swarm.ReplicatedService{
		Replicas: &replicas,
	}
	response, err := scaler.Client.ServiceUpdate(ctx, service.ID, service.Meta.Version, service.Spec, types.ServiceUpdateOptions{})

	if err != nil {
		log.Error(err.Error())
		return err
	}

	if len(response.Warnings) > 0 {
		log.Warnf("received scaling up service %s: %v", name, response.Warnings)
	}

	return nil
}

func (scaler *DockerSwarmScaler) IsUp(name string) bool {
	ctx := context.Background()
	service, err := scaler.GetServiceByName(name, ctx)

	if err != nil {
		log.Error(err.Error())
		return false
	}

	return service.ServiceStatus.DesiredTasks > 0 && (service.ServiceStatus.DesiredTasks == service.ServiceStatus.RunningTasks)
}

func (scaler *DockerSwarmScaler) GetServiceByName(name string, ctx context.Context) (*swarm.Service, error) {
	opts := types.ServiceListOptions{
		Filters: filters.NewArgs(),
		Status:  true,
	}
	opts.Filters.Add("name", name)

	services, err := scaler.Client.ServiceList(ctx, opts)

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
