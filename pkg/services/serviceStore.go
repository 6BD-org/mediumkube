package services

import "mediumkube/pkg/common"

var k8sService KubernetesService
var libvirtService LibvirtService
var dnsService *DNSService
var meshService *MeshService

func InitK8sService(ks KubernetesService) {
	k8sService = ks
}

func GetK8sService() KubernetesService {
	return k8sService
}

func InitLibvritService(ls LibvirtService) {
	libvirtService = ls
}

func InitMeshService(ms MeshService) {
	meshService = &ms
}

func GetLibvirtService() LibvirtService {
	return libvirtService
}

func GetMeshService() *MeshService {
	return meshService
}

func GetNodeManager(backend string) NodeManager {
	switch backend {
	case "libvirt":
		return libvirtService
	default:
		panic("Unsupported type")
	}
}

func GetDNSService(config *common.OverallConfig) *DNSService {
	if dnsService == nil {
		dnsService = NewDNSService(config)
	}
	return dnsService
}
