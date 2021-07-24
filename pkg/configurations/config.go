package configurations

import (
	"log"
	"mediumkube/pkg/common"
	"mediumkube/pkg/utils"

	"gopkg.in/yaml.v3"
)

const (
	configDir = "/etc/mediumkube/config.yaml"
)

var overallConfig *common.OverallConfig = nil

func LoadConfigFromFile(configPath string) *common.OverallConfig {
	configStr := utils.ReadByte(configPath)
	_config := &common.OverallConfig{}
	err := yaml.Unmarshal(configStr, _config)
	utils.CheckErr(err)
	return _config
}

// InitConfig initialize configuration context
func InitConfig() {
	log.Println("Using configuration file: ", configDir)
	overallConfig = LoadConfigFromFile(configDir)
}

func EnsureFlavorExistance(config *common.OverallConfig) {
	if len(config.Flavors) < 1 {
		panic("Must specify at least one flavor")
	}
}

// Config Get config
func Config() *common.OverallConfig {
	if overallConfig == nil {
		InitConfig()
	}
	EnsureFlavorExistance(overallConfig)
	return overallConfig
}
