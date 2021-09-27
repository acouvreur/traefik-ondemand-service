package scaler

var onereplicas = uint64(1)
var zeroreplicas = uint64(0)

type Scaler interface {
	ScaleUp(name string) error
	ScaleDown(name string) error
	IsUp(name string) bool
}
