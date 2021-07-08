package services

import (
	"context"
	"encoding/json"
	"mediumkube/pkg/common"
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/etcd"
	"mediumkube/pkg/models"

	"go.etcd.io/etcd/client/v2"
	"k8s.io/klog/v2"
)

type MeshService struct {
	config *common.OverallConfig
}

func (m *MeshService) ListDomains() ([]models.Domain, error) {
	res := make([]models.Domain, 0)
	prefix := m.config.Overlay.DomainEtcdPrefix
	kpi := client.NewKeysAPI(etcd.NewClientOrDie())
	resp, err := kpi.Get(context.TODO(), prefix, nil)
	if err != nil {
		return res, nil
	}

	for _, n := range resp.Node.Nodes {
		ds := make([]models.Domain, 0)
		err = json.Unmarshal([]byte(n.Value), &ds)
		if err != nil {
			klog.Error(err)
		}
		for _, d := range ds {
			res = append(res, d)
		}
	}
	return res, nil
}

func init() {
	InitMeshService(MeshService{config: configurations.Config()})
}
