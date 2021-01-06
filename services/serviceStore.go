package services

var multipassService MultipassService
var k8sService KubernetesService
var libvirtService LibvirtService

func InitMultipassService(ms MultipassService) {
	multipassService = ms
}

func GetMultipassService() MultipassService {
	return multipassService
}

func InitK8sService(ks KubernetesService) {
	k8sService = ks
}

func GetK8sService() KubernetesService {
	return k8sService
}

func InitLibvritService(ls LibvirtService) {
	libvirtService = ls
}

func GetLibvirtService() LibvirtService {
	return libvirtService
}

func GetNodeManager(backend string) NodeManager {
	switch backend {
	case "multipass":
		return multipassService
	case "libvirt":
		return libvirtService
	default:
		panic("Unsupported type")
	}
}
