package configurations

import (
	"log"
	"mediumkube/common"
	"mediumkube/utils"

	"gopkg.in/yaml.v3"
)

var overallConfig *common.OverallConfig = &common.OverallConfig{}

// InitConfig initialize configuration context
func InitConfig(configDir string) {
	log.Println("Using configuration file: ", configDir)
	configStr := utils.ReadByte(configDir)
	err := yaml.Unmarshal(configStr, overallConfig)
	utils.CheckErr(err)
}

// Config Get config
func Config() *common.OverallConfig {
	return overallConfig
}
