package common

import "path"

// VolumeMount mount host dir int virtual machine dir
type VolumeMount struct {
	Host string `yaml:"host"`
	VM   string `yaml:"vm"`
}

// Network config for node in libvirt mode
type Network struct {
	Name string `yaml:"name"`
	IP   string `yaml:"ip"`
}

// NodeConfig Config for each noed. "node" field  in config yaml
type NodeConfig struct {
	Name string `yaml:"name"`
	CPU  string `yaml:"cpu"`
	MEM  string `yaml:"mem"`
	DISK string `yaml:"disk"`
}

// Arg Argument
type Arg struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// KubeInit Configuration for kubeadm init
type KubeInit struct {
	Args []Arg `yaml:"args"`
}

// Bridge that establish connections between vms
type Bridge struct {
	Name      string `yaml:"name"`
	Alias     string `yaml:"alias"`
	Inet      string `yaml:"inet"`
	Broadcast string `yaml:"broadcast"`
	Host      string `yaml:"host"` // NIC on host machine
}

// OverallConfig Over-all configs
/*
This config is in same file with Template config, but when parsed,
the template part will be ignored.
*/
type OverallConfig struct {
	HTTPSProxy      string       `yaml:"https-proxy,omitempty"`
	HTTPProxy       string       `yaml:"http-proxy,omitempty"`
	Backend         string       `yaml:"backend"`
	Bridge          Bridge       `yaml:"bridge"`
	NodeConfig      []NodeConfig `yaml:"nodes"`
	Image           string       `yaml:"image"`
	CloudInit       string       `yaml:"cloud-init"`
	KubeInit        KubeInit     `yaml:"kube-init"`
	TmpDir          string       `yaml:"tmp_dir"`            // Directory to store temp files generated by medium kube
	VMKubeConfigDir string       `yaml:"vm_kube_config_dir"` // path of kubeconfig in virtual machine
	PubKeyDir       string       `yaml:"pub-key-dir"`
	PrivKeyDir      string       `yaml:"priv-key-dir"`
	HostPubKeyDir   string       `yaml:"host-pub-key-dir"`
	HostPrivKeyDir  string       `yaml:"host-priv-key-dir"` // Used for ssh execute
}

type KVMDomainConfig struct {
}

func (config OverallConfig) LeaseFile() string {
	return path.Join(config.TmpDir, "dnsmasq.lease")
}
