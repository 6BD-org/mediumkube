package common

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

// NodeConfig Config for each noed. "node" field  in config yaml
type NodeConfig struct {
	CPU  string `yaml:"cpu"`
	MEM  string `yaml:"gpu"`
	DISK string `yaml:"disk"`
}

// OverallConfig Over-all configs
/*
This config is in same file with Template config, but when parsed,
the template part will be ignored.
*/
type OverallConfig struct {
	NodeConfig NodeConfig `yaml:"node"`
	NodeNum    int        `yaml:"node-num"`
}
