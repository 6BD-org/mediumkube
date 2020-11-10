package k8s

import (
	"io"
	"mediumkube/common"
	"mediumkube/utils"
	"os"
	"path/filepath"

	go_yaml "gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// KubeConfigPath Get path of kube config on host machine
func KubeConfigPath(config common.OverallConfig) string {
	return filepath.Join(config.TmpDir, ".kube/config")
}

func ParseResources(path string) map[string]interface{} {
	resourceMapping := NewResourceMapping()
	reader, err := os.Open(path)
	utils.CheckErr(err)
	splittedReader := yaml.NewDocumentDecoder(reader)

	for {
		buf := make([]byte, 5*1024*1024)
		kindAware := KindAware{}
		_, err := splittedReader.Read(buf)
		if err != nil && err != io.EOF {
			utils.CheckErr(err)
		}
		go_yaml.Unmarshal(buf, &kindAware)
		go_yaml.Unmarshal(buf, resourceMapping[kindAware.Kind])
		if err == io.EOF {
			break
		}
	}

	return resourceMapping
}
