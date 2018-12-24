package utils

import (
	"fmt"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"

	redisclientset "github.com/zdq0394/redis-cluster-operator/pkg/client/clientset/versioned"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// LoadKubernetesConfig loads kubernetes configuration based on flags.
func LoadKubernetesConfig(deployment bool) (*rest.Config, error) {
	var cfg *rest.Config
	// If devel mode then use configuration flag path.
	kubeConfig := "/etc/kubernetes/kubeconfig"
	if deployment {
		config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			return nil, fmt.Errorf("could not load configuration: %s", err)
		}
		cfg = config
	} else {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("error loading kubernetes configuration inside cluster, check app is running outside kubernetes cluster or run in development mode: %s", err)
		}
		cfg = config
	}

	return cfg, nil
}

// CreateKubernetesClients create the clients to connect to kubernetes
func CreateKubernetesClients(development bool) (kubernetes.Interface, redisclientset.Interface, apiextensionsclientset.Interface, error) {
	config, err := LoadKubernetesConfig(development)
	if err != nil {
		return nil, nil, nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, err
	}
	customClientset, err := redisclientset.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, err
	}

	aeClientset, err := apiextensionsclientset.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, err
	}

	return clientset, customClientset, aeClientset, nil
}
