package mgrpc

import (
	context "context"
	"fmt"
	"mediumkube/pkg/common"
	"mediumkube/pkg/common/event"
	"mediumkube/pkg/dlock"
	"mediumkube/pkg/services"
	"strings"

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
	manager := services.GetDomainManager(s.config.Backend)
	domainList, err := manager.List()
	if err != nil {
		return nil, err
	}
	resp := DomainListResp{}
	resp.state = protoimpl.MessageState{}
	resp.Domains = MarshalList(domainList)
	return &resp, nil
}

func (s *MediumKubeServer) DeployDomain(param *DomainCreationParam, stream DomainSercice_DeployDomainServer) error {
	manager := services.GetDomainManager(s.config.Backend)
	var err error
	dlock.NewEtcdLockManager(s.config).DoWithLock(domainCreationLockType, 5*60*1000*1000*1000, func() {
		ms := services.GetMeshService()
		existingDomains, err := ms.ListDomains()
		if err != nil {
			return
		}
		sink := func(b []byte) error {
			resp := DomainCreationResp{
				Content: b,
			}
			return stream.Send(&resp)
		}
		ncs := make([]common.NodeConfig, 0)
		for _, config := range param.Config {
			for _, ed := range existingDomains {
				if strings.ToLower(ed.Name) == strings.ToLower(config.Name) {
					err = fmt.Errorf("Domain already exists")
					return
				}
			}

			ncs = append(ncs, common.NodeConfig{
				CPU:       config.Cpu,
				MEM:       config.Memory,
				DISK:      config.Disk,
				Name:      config.Name,
				CloudInit: config.CloudInit,
			})
		}
		manager.Deploy(ncs, s.config.Image, sink)
	}, func() {
		err = fmt.Errorf("Failed to acquire lock, giving up deployment")
	})
	if err != nil {
		return err
	}

	event.GetEventBus().DomainUpdate <- event.DomainEvent{}
	return nil
}

func (s *MediumKubeServer) DeleteDomains(ctx context.Context, param *DomainDeletionParam) (*DomainDeletionResp, error) {
	manager := services.GetDomainManager(s.config.Backend)
	for _, domain := range param.Names {
		manager.Purge(domain)
	}
	return &DomainDeletionResp{}, nil
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
