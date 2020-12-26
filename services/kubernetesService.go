package services

import (
	"context"
	"log"
	"mediumkube/common"
	"mediumkube/configurations"
	"mediumkube/k8s"
	"mediumkube/utils"

	appsV1 "k8s.io/api/apps/v1"
	appsV1beta1 "k8s.io/api/apps/v1beta1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	rbacV1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// KubernetesService service to interact with kubernetes
type KubernetesService struct {
	OverallConfig *common.OverallConfig
}

// Client initialize a k8s client from overall config
func (service KubernetesService) Client() *kubernetes.Clientset {
	// kubeconfig := flag.String("kubeconfig", k8s.KubeConfigPath(service.OverallConfig), "Path to kubeconfig file")
	log.Println("Using kube config file: ", service.OverallConfig, k8s.KubeConfigPath(service.OverallConfig))
	config, _ := clientcmd.BuildConfigFromFlags("", k8s.KubeConfigPath(service.OverallConfig))
	clientset, err := kubernetes.NewForConfig(config)
	utils.CheckErr(err)
	return clientset
}

func parse(path string) {
	// file, err = os.Open(path)
	// utils.CheckErr(err)
	// decoder := yaml.NewYAMLOrJSONDecoder(file, 100)
}

// Apply Install a yaml resource to cluster
func (service KubernetesService) Apply(file string) {
	client := service.Client()
	ch := make(chan interface{})
	go k8s.ParseResources(file, ch)
	ctx := context.TODO()
	for v := range ch {
		switch v.(type) {
		case *v1beta1.PodSecurityPolicy:
			field, ok := v.(*v1beta1.PodSecurityPolicy)
			if ok {
				log.Println("Installing K8s resource PodSecurityPolicies: ", field.Name)
				client.ExtensionsV1beta1().PodSecurityPolicies().Create(ctx, field, v1.CreateOptions{})
			}
		case *rbacV1.ClusterRole:
			field, ok := v.(*rbacV1.ClusterRole)
			if ok {
				log.Println("Installing K8s resource ClusterRoles: ", field.Name)
				client.RbacV1().ClusterRoles().Create(ctx, field, v1.CreateOptions{})
			}
		case *rbacV1.ClusterRoleBinding:
			field, ok := v.(*rbacV1.ClusterRoleBinding)
			if ok {
				log.Println("Installing K8s resource ClusterRoleBindings: ", field.Name)
				client.RbacV1().ClusterRoleBindings().Create(ctx, field, v1.CreateOptions{})
			}
		case *coreV1.ServiceAccount:
			field, ok := v.(*coreV1.ServiceAccount)
			if ok {
				log.Println("Installing K8s resource ServiceAccounts: ", field.Name)
				client.CoreV1().ServiceAccounts(field.Namespace).Create(ctx, field, v1.CreateOptions{})
			}
		case *coreV1.ConfigMap:
			field, ok := v.(*coreV1.ConfigMap)
			if ok {
				log.Println("Installing K8s resource ConfigMaps: ", field.Name)
				client.CoreV1().ConfigMaps(field.Namespace).Create(ctx, field, v1.CreateOptions{})
			}
		case *appsV1.DaemonSet:
			field, ok := v.(*appsV1.DaemonSet)
			if ok {
				log.Println("Installing K8s resource DaemonSets: ", field.Name)
				client.AppsV1().DaemonSets(field.Namespace).Create(ctx, field, v1.CreateOptions{})
			}
		case *appsV1.StatefulSet:
			field, ok := v.(*appsV1.StatefulSet)
			if ok {
				log.Println("Installing K8s resource StatefulSets: ", field.Name)
				client.AppsV1().StatefulSets(field.Namespace).Create(ctx, field, v1.CreateOptions{})
			}
		case *coreV1.Service:
			field, ok := v.(*coreV1.Service)
			if ok {
				log.Println("Installing K8s resource Service: ", field.Name)
				client.CoreV1().Services(field.Namespace).Create(ctx, field, v1.CreateOptions{})
			}
		case *appsV1beta1.Deployment:
			field, ok := v.(*appsV1beta1.Deployment)
			if ok {
				log.Println("Installing K8s resource Deployment: ", field.Name)
				client.AppsV1beta1().Deployments(field.Namespace).Create(ctx, field, v1.CreateOptions{})
			}
		}
	}
}

func init() {
	InitK8sService(KubernetesService{
		OverallConfig: configurations.Config(),
	})
}
