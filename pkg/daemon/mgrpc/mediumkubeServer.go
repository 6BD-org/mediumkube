package mgrpc

import (
	context "context"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type MediumKubeServer struct {
	UnimplementedDomainSerciceServer
}

func (s *MediumKubeServer) ListDomains(context.Context, *EmptyParam) (*DomainListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDomains not implemented")
}
func (s *MediumKubeServer) DeployDomain(context.Context, *DomainCreationParam) (*DomainCreationResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeployDomain not implemented")
}

func init() {
	var _ DomainSerciceServer = (*MediumKubeServer)(nil)
}
