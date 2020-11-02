package services

var multipassService MultipassService
var k8sService KubernetesService

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
