package k8s

import (
	"bytes"
	"io"
	"mediumkube/common"
	"mediumkube/utils"
	"os"
	"path/filepath"

	go_yaml "gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// KubeConfigPath Get path of kube config on host machine
func KubeConfigPath(config *common.OverallConfig) string {
	return filepath.Join(config.TmpDir, ".kube/config")
}

// ParseResources Parse Yaml to k8s resource
func ParseResources(path string) map[string]interface{} {
	resourceMapping := NewResourceMapping()
	reader, err := os.Open(path)
	utils.CheckErr(err)
	splittedReader := yaml.NewDocumentDecoder(reader)

	for {
		buf := make([]byte, 5*1024*1024)

		kindIdentifier := KindIdentifier{}
		size, err := splittedReader.Read(buf)
		if size == 0 {
			break
		}
		if err != nil && err != io.EOF {
			utils.CheckErr(err)
		}
		err = go_yaml.Unmarshal(buf[:size], &kindIdentifier)
		utils.CheckErr(err)
		decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(buf[:size]), size)
		err = decoder.Decode(resourceMapping[kindIdentifier.Kind])
		utils.CheckErr(err)
		if err == io.EOF {
			break
		}
	}

	return resourceMapping
}
