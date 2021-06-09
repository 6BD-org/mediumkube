package mesh

type PeerLease struct {
	Master    bool   `yaml:"master"`
	Cidr      string `yaml:"cidr"`
	PublicIP  string `yaml:"publicIP"`
	Timestamp int64  `yaml:"timestamp"`
	TTL       int64  `yaml:"ttl"`
}
