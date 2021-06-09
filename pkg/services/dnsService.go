package services

import (
	"mediumkube/pkg/common"
	"mediumkube/pkg/network"
)

type DNSService struct {
	config *common.OverallConfig
}

func NewDNSService(config *common.OverallConfig) *DNSService {
	return &DNSService{
		config: config,
	}
}

func (dns *DNSService) isMeshEnabled() bool {
	return dns.config.Overlay.Enabled
}

func (dns *DNSService) Resolve(hostname string) (string, bool) {
	if dns.isMeshEnabled() {
		return network.Resolve(dns.config.LeaseFile(), hostname)
	} else {
		// ETCD not supported yet
		// return network.Resolve(dns.config.LeaseFile(), hostname)
		return network.ResolveOverlay(dns.config.Overlay, hostname)
	}
}
