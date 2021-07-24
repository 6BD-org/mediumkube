package services

import (
	"context"
	"encoding/json"
	"mediumkube/pkg/common"
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/etcd"
	"mediumkube/pkg/models"

	"go.etcd.io/etcd/client/v2"
	clientv2 "go.etcd.io/etcd/client/v2"
	"k8s.io/klog/v2"
)

type MeshService struct {
	config *common.OverallConfig
}

func (m *MeshService) ListLeases() ([]models.PeerLease, error) {
	res := make([]models.PeerLease, 0)
	kpi := clientv2.NewKeysAPI(etcd.NewClientOrDie())
	resp, err := kpi.Get(context.TODO(), m.config.Overlay.LeaseEtcdPrefix, nil)
	if err != nil {
		klog.Error(err)
		return []models.PeerLease{}, err
	}

	for _, node := range resp.Node.Nodes {
		payload := models.PeerLease{}
		if len(node.Value) == 0 {
			continue
		}
		err = json.Unmarshal([]byte(node.Value), &payload)
		if err != nil {

			klog.Errorf("Fail to marshal payload: %v, err: %v", node.Value, err)
			continue
		}
		res = append(res, payload)
	}
	return res, nil
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
