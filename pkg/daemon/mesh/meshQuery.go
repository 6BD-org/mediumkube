package mesh

import (
	"context"
	"encoding/json"
	"mediumkube/pkg/common"
	etcd "mediumkube/pkg/etcd"
	"mediumkube/pkg/models"

	"go.etcd.io/etcd/client/v2"
	"k8s.io/klog/v2"
)

func ClusterLeases(config *common.OverallConfig) ([]models.PeerLease, error) {
	return pullLease(config)
}

func ListDomains(config *common.OverallConfig) ([]models.Domain, error) {
	res := make([]models.Domain, 0)
	prefix := config.Overlay.DomainEtcdPrefix
	if etcdClient == nil {
		etcdClient = etcd.NewClientOrDie()
	}
	kpi := client.NewKeysAPI(etcdClient)
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
