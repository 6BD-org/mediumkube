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

// InitConfig initialize configuration context
func InitConfig() {
	log.Println("Using configuration file: ", configDir)
	configStr := utils.ReadByte(configDir)
	err := yaml.Unmarshal(configStr, overallConfig)
	utils.CheckErr(err)
}

// Config Get config
func Config() *common.OverallConfig {
	if overallConfig == nil {
		InitConfig()
	}
	return overallConfig
}
