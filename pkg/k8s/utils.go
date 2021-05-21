package k8s

import (
	"bytes"
	"io"
	"mediumkube/pkg/common"
	"mediumkube/pkg/utils"
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
// For compositional yaml files, each component is decoded and written into channel
func ParseResources(path string, ch chan interface{}) {
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

		ch <- resourceMapping[kindIdentifier.Kind]
		if err == io.EOF {
			break
		}
	}
	close(ch)

}
