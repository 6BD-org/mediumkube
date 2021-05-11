package flannel

import (
	"encoding/json"
	"mediumkube/common"
)

// FlannelConfig in etcd
type FlannelConfig struct {
	Network string  `json:"Network"`
	Backend Backend `json:"Backend"`
}

// Flannel Backend
type Backend struct {
	Type string `json:"Type"`
}

// NewConfig creates config from overall config
func NewConfig(config *common.OverallConfig) *FlannelConfig {
	return &FlannelConfig{
		Network: config.Overlay.Flannel.Network,
		Backend: Backend{
			Type: config.Overlay.Flannel.Backend,
		},
	}
}

// ToStr converts config object to string
func (fc *FlannelConfig) ToStr() string {
	bts, err := json.Marshal(fc)
	if err != nil {
		panic(err)
	}
	return string(bts)
}
