package services

import (
	"flag"
	"mediumkube/common"
	"mediumkube/configurations"
	"mediumkube/k8s"
	"mediumkube/utils"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// KubernetesService service to interact with kubernetes
type KubernetesService struct {
	OverallConfig common.OverallConfig
}

// Client initialize a k8s client from overall config
func (service *KubernetesService) Client() *kubernetes.Clientset {
	kubeconfig := flag.String("kubeconfig", k8s.KubeConfigPath(service.OverallConfig), "Path to kubeconfig file")
	config, _ := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	clientset, err := kubernetes.NewForConfig(config)
	utils.CheckErr(err)
	return clientset
}

func parse(path string) {
	// file, err = os.Open(path)
	// utils.CheckErr(err)
	// decoder := yaml.NewYAMLOrJSONDecoder(file, 100)
}

func (service *KubernetesService) Apply(file string) {
	// client := service.Client()
	// client.AppsV1().DaemonSets().Create()
}

func init() {
	InitK8sService(KubernetesService{
		OverallConfig: configurations.Config(),
	})
}
