package scaler

import "github.com/docker/docker/client"

type Scaler interface {
	ScaleUp(client *client.Client, name string, replicas *uint64)
	ScaleDown(client *client.Client, name string)
	IsUp(client *client.Client, name string) bool
}
