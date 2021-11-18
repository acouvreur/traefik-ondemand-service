package scaler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Delimiter is used to split name into kind,namespace,name,replicacount
const Delimiter = "_"

type Config struct {
	Kind      string // deployment or statefulset
	Namespace string
	Name      string
	Replicas  int
}

type Workload interface {
	GetScale(ctx context.Context, workloadName string, options metav1.GetOptions) (*autoscalingv1.Scale, error)
	UpdateScale(ctx context.Context, workloadName string, scale *autoscalingv1.Scale, opts metav1.UpdateOptions) (*autoscalingv1.Scale, error)
}

func convertName(name string) (*Config, error) {
	// name format kind_namespace_name_replicas
	s := strings.Split(name, Delimiter)
	if len(s) < 4 {
		return nil, errors.New("invalid name should be: kind" + Delimiter + "namespace" + Delimiter + "name" + Delimiter + "replicas")
	}
	replicas, err := strconv.Atoi(s[3])
	if err != nil {
		return nil, err
	}

	return &Config{
		Kind:      s[0],
		Namespace: s[1],
		Name:      s[2],
		Replicas:  replicas,
	}, nil
}

type KubernetesScaler struct {
	Client *kubernetes.Clientset
}

func NewKubernetesScaler(client *kubernetes.Clientset) *KubernetesScaler {
	return &KubernetesScaler{
		Client: client,
	}
}

func (scaler *KubernetesScaler) ScaleUp(name string) error {
	config, err := convertName(name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Infof("Scaling up %s %s in namespace %s to %d", config.Kind, config.Name, config.Namespace, config.Replicas)
	ctx := context.Background()

	var workload Workload

	switch config.Kind {
	case "deployment":
		workload = scaler.Client.AppsV1().Deployments(config.Namespace)
	case "statefulset":
		workload = scaler.Client.AppsV1().StatefulSets(config.Namespace)
	default:
		return fmt.Errorf("unsupported kind %s", config.Kind)
	}

	s, err := workload.GetScale(ctx, config.Name, metav1.GetOptions{})
	if err != nil {
		log.Error(err.Error())
		return err
	}

	sc := *s
	if sc.Spec.Replicas == 0 {
		sc.Spec.Replicas = int32(config.Replicas)
	} else {
		log.Infof("Replicas for %s %s in namespace %s are already: %d", config.Kind, config.Name, config.Namespace, sc.Spec.Replicas)
		return nil
	}

	_, err = workload.UpdateScale(ctx, config.Name, &sc, metav1.UpdateOptions{})

	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (scaler *KubernetesScaler) ScaleDown(name string) error {
	config, err := convertName(name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Infof("Scaling down %s %s in namespace %s to 0", config.Kind, config.Name, config.Namespace)
	ctx := context.Background()

	var workload Workload

	switch config.Kind {
	case "deployment":
		workload = scaler.Client.AppsV1().Deployments(config.Namespace)
	case "statefulset":
		workload = scaler.Client.AppsV1().StatefulSets(config.Namespace)
	default:
		return fmt.Errorf("unsupported kind %s", config.Kind)
	}

	s, err := workload.GetScale(ctx, config.Name, metav1.GetOptions{})
	if err != nil {
		log.Error(err.Error())
		return err
	}

	sc := *s
	if sc.Spec.Replicas != 0 {
		sc.Spec.Replicas = 0
	} else {
		log.Infof("Replicas for %s %s in namespace %s are already: 0", config.Kind, config.Name, config.Namespace)
		return nil
	}

	_, err = workload.UpdateScale(ctx, config.Name, &sc, metav1.UpdateOptions{})

	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (scaler *KubernetesScaler) IsUp(name string) bool {
	ctx := context.Background()

	config, err := convertName(name)
	if err != nil {
		log.Error(err.Error())
		return false
	}

	switch config.Kind {
	case "deployment":
		d, err := scaler.Client.AppsV1().Deployments(config.Namespace).
			Get(ctx, config.Name, metav1.GetOptions{})
		if err != nil {
			log.Error(err.Error())
			return false
		}
		log.Infof("Status for %s %s in namespace %s is: AvailableReplicas %d, ReadyReplicas: %d ", config.Kind, config.Name, config.Namespace, d.Status.AvailableReplicas, d.Status.ReadyReplicas)

		if d.Status.ReadyReplicas > 0 {
			return true
		}
	case "statefulset":
		d, err := scaler.Client.AppsV1().StatefulSets(config.Namespace).
			Get(ctx, config.Name, metav1.GetOptions{})
		if err != nil {
			log.Error(err.Error())
			return false
		}
		log.Infof("Status for %s %s in namespace %s is: AvailableReplicas %d, ReadyReplicas: %d ", config.Kind, config.Name, config.Namespace, d.Status.AvailableReplicas, d.Status.ReadyReplicas)

		if d.Status.ReadyReplicas > 0 {
			return true
		}

	default:
		log.Error(fmt.Errorf("unsupported kind %s", config.Kind))
		return false
	}

	return false
}
