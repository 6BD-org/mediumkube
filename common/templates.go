package common

import (
	"mediumkube/utils"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// TemplateConfig this is parsed from config yaml
type TemplateConfig struct {
	HTTPSProxy    string `yaml:"https-proxy,omitempty"`
	HTTPProxy     string `yaml:"http-proxy,omitempty"`
	PubKeyDir     string `yaml:"pub-key-dir"`
	PrivKeyDir    string `yaml:"priv-key-dir"`
	HostPubKeyDir string `yaml:"host-pub-key-dir"`

	PubKey     string
	PrivKey    string
	HostPubKey string
}

// LoadFile load the content of of file into template
func (TC TemplateConfig) LoadFile(path string) string {
	return utils.ReadStr(path)
}

// Node Config templates

// MemoryInBytes get mem in bytes
func (nc NodeConfig) MemoryInBytes() int64 {
	s := strings.ReplaceAll(nc.MEM, "G", "")
	gb, err := strconv.ParseFloat(s, 64)
	utils.CheckErr(err)
	return int64(gb * (1000000000))
}

func (nc NodeConfig) MemoryInMiB() int64 {
	s := strings.ReplaceAll(nc.MEM, "G", "")
	gb, err := strconv.ParseFloat(s, 64)
	utils.CheckErr(err)
	return int64(gb * (1000))
}

// UUID generate uuid for node
func (nc NodeConfig) UUID() string {
	return uuid.New().String()
}
