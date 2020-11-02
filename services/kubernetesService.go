package services

import (
	"mediumkube/common"
	"mediumkube/configurations"
	"path/filepath"
)

// KubernetesService service to interact with kubernetes
type KubernetesService struct {
	OverallConfig common.OverallConfig
}

func configPath(config common.OverallConfig) string {
	return filepath.Join(config.TmpDir, ".kube/config")
}

func Apply(file string) {
}

func init() {
	InitK8sService(KubernetesService{
		OverallConfig: configurations.Config(),
	})
}
