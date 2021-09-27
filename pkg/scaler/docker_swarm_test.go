package scaler

import (
	"context"
	"reflect"
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

func TestDockerSwarmScaler_ScaleUp(t *testing.T) {
	type fields struct {
		Client client.ServiceAPIClient
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scaler := &DockerSwarmScaler{
				Client: tt.fields.Client,
			}
			scaler.ScaleUp(tt.args.name)
		})
	}
}

func TestDockerSwarmScaler_ScaleDown(t *testing.T) {
	type fields struct {
		Client client.ServiceAPIClient
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scaler := &DockerSwarmScaler{
				Client: tt.fields.Client,
			}
			scaler.ScaleDown(tt.args.name)
		})
	}
}

func TestDockerSwarmScaler_IsUp(t *testing.T) {
	type fields struct {
		Client client.ServiceAPIClient
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scaler := &DockerSwarmScaler{
				Client: tt.fields.Client,
			}
			if got := scaler.IsUp(tt.args.name); got != tt.want {
				t.Errorf("DockerSwarmScaler.IsUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDockerSwarmScaler_GetServiceByName(t *testing.T) {
	type fields struct {
		Client client.ServiceAPIClient
	}
	type args struct {
		name string
		ctx  context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *swarm.Service
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scaler := &DockerSwarmScaler{
				Client: tt.fields.Client,
			}
			got, err := scaler.GetServiceByName(tt.args.name, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("DockerSwarmScaler.GetServiceByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DockerSwarmScaler.GetServiceByName() = %v, want %v", got, tt.want)
			}
		})
	}
}
