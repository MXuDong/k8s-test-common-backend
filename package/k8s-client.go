package client

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

const (
	EnvIsInCluster = "IS_IN_CLUSTER"

	InCluster = "true"
)

//InitClient will return kubectl client, the return is client set, is out side of cluster and init error
func InitClient(configPath string) (*kubernetes.Clientset, *rest.Config, bool, error) {
	isInCluster := IsInKubernetes()
	if isInCluster {
		clientSet, config, err := insideMode()
		return clientSet, config, isInCluster, err
	} else {
		clientSet, config, err := outsideMode(configPath)
		return clientSet, config, isInCluster, err
	}
}

// IsInKubernetes will return the application is run in kubernetes, it get value from env : EnvIsInCluster
func IsInKubernetes() bool {
	return os.Getenv(EnvIsInCluster) == InCluster
}

//outsideMode will init kubernetes client out side of cluster
func outsideMode(configPath string) (*kubernetes.Clientset, *rest.Config, error) {

	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return clientSet, config, nil
}

// insideMode will init kubernetes client in side of cluster
func insideMode() (*kubernetes.Clientset, *rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return clientSet, config, nil
}
