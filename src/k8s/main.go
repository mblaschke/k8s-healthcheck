package k8s

import (
	"k8s-healthcheck/src/logger"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	coreV1 "k8s.io/api/core/v1"
)

type KubernetesBase struct {
	clientset *kubernetes.Clientset

	Logger *logger.DaemonLogger
	KubeContext string
	KubeConfig string
}

type Kubernetes struct {
	KubernetesBase
}


// Create cached kubernetes client
func (k *KubernetesBase) Client() (clientset *kubernetes.Clientset) {
	var err error
	var config *rest.Config

	if k.clientset == nil {
		if k.KubeConfig != "" {
			// KUBECONFIG
			config, err = buildConfigFromFlags(k.KubeContext, k.KubeConfig)
			if err != nil {
				panic(err.Error())
			}
		} else {
			// K8S in cluster
			config, err = rest.InClusterConfig()
			if err != nil {
				panic(err.Error())
			}
		}

		k.clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	}

	return k.clientset
}

func buildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}


func (k *KubernetesBase) NodeList() (*coreV1.NodeList, error) {
	options := v1.ListOptions{}
	return k.Client().CoreV1().Nodes().List(options)
}
