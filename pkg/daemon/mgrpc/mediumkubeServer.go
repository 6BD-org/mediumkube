package mgrpc

import (
	context "context"
	"mediumkube/pkg/common"
	"mediumkube/pkg/services"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/runtime/protoimpl"
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
	return nil, status.Errorf(codes.Unimplemented, "method DeployDomain not implemented")
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
