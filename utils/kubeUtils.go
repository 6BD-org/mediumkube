package utils

import (
	"mediumkube/common"
	"path/filepath"
)

// KubeConfigPath Get path of kube config on host machine
func KubeConfigPath(config common.OverallConfig) string {
	return filepath.Join(config.TmpDir, ".kube/config")
}
