package mgrpc

import (
	context "context"
	"mediumkube/pkg/common"
	"mediumkube/pkg/common/event"
	"mediumkube/pkg/dlock"
	"mediumkube/pkg/services"

	"google.golang.org/protobuf/runtime/protoimpl"
)

const (
	domainCreationLockType = "LOCK_DOMAIN_CREATION"
)

type MediumKubeServer struct {
	UnimplementedDomainSerciceServer
	config *common.OverallConfig
}

func (s *MediumKubeServer) ListDomains(context.Context, *EmptyParam) (*DomainListResp, error) {
	manager := services.GetNodeManager(s.config.Backend)
	domainList, err := manager.List()
	if err != nil {
		return nil, err
	}
	resp := DomainListResp{}
	resp.state = protoimpl.MessageState{}
	resp.Domains = MarshalList(domainList)
	return &resp, nil
}

func (s *MediumKubeServer) DeployDomain(context.Context, *DomainCreationParam) (*DomainCreationResp, error) {
	manager := services.GetNodeManager(s.config.Backend)
	dlock.NewEtcdLockManager(s.config).DoWithLock(domainCreationLockType, 5*60*1000*1000*1000, func() {
		manager.Deploy(make([]common.NodeConfig, 0), "", "")
	}, func() {})
	event.GetEventBus().DomainUpdate <- event.DomainEvent{}
	return nil, nil
}

func NewServer(config *common.OverallConfig) *MediumKubeServer {
	s := MediumKubeServer{
		config: config,
	}
	return &s
}

func init() {
	var _ DomainSerciceServer = (*MediumKubeServer)(nil)
}
