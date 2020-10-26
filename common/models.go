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
