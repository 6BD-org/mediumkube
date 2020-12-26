package k8s

import (
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	v1 "k8s.io/api/rbac/v1"
)

type KindIdentifier struct {
	Kind string `yaml:"kind"`
}

// New create a new resource mapping
func NewResourceMapping() map[string]interface{} {
	var resourceType = make(map[string]interface{})
	resourceType["PodSecurityPolicy"] = &v1beta1.PodSecurityPolicy{}
	resourceType["ClusterRole"] = &v1.ClusterRole{}
	resourceType["ClusterRoleBinding"] = &v1.ClusterRoleBinding{}
	resourceType["ServiceAccount"] = &coreV1.ServiceAccount{}
	resourceType["ConfigMap"] = &coreV1.ConfigMap{}
	resourceType["DaemonSet"] = &appsV1.DaemonSet{}
	resourceType["StatefulSet"] = &appsV1.StatefulSet{}
	resourceType["Deployment"] = &v1beta1.Deployment{}
	resourceType["Service"] = &coreV1.Service{}

	return resourceType
}
