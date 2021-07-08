package mesh

import (
	"context"
	"encoding/json"
	"mediumkube/pkg/common"
	"mediumkube/pkg/common/flannel"
	"mediumkube/pkg/configurations"
	etcd "mediumkube/pkg/etcd"
	"mediumkube/pkg/models"
	"mediumkube/pkg/network"
	"mediumkube/pkg/services"
	"strings"
	"time"

	clientv2 "go.etcd.io/etcd/client/v2"
	"k8s.io/klog/v2"
)

const (
	leaseTTL = 60 // 6 seconds
)

func _nsKey(config *common.OverallConfig, k string) string {
	prefix := config.Overlay.DNSEtcdPrefix
	return prefix + "/" + k
}

func _leaseKey(config *common.OverallConfig) string {
	return config.Overlay.LeaseEtcdPrefix + "/" + strings.Replace(config.Overlay.Cidr, "/", "-", -1)
}

func _domainKey(config *common.OverallConfig) string {
	return config.Overlay.DomainEtcdPrefix + "/" + strings.Replace(config.Overlay.Cidr, "/", "-", -1)
}

func doSync(config *common.OverallConfig) {
	locals := network.ListNSPairs(config.DomainNSFile())
	remotes := network.ListETCDNsPairs(etcd.NewClientOrDie(), config.Overlay.DNSEtcdPrefix, config.Overlay.Cidr)
	in, out := network.SyncDNSLease(locals, remotes)

	kpi := clientv2.NewKeysAPI(etcd.NewClientOrDie())

	for _, pairIn := range in {
		// PUT etcd
		val, marshalErr := json.Marshal(&pairIn)
		if marshalErr != nil {
			klog.Error(marshalErr)
		} else {
			_, err := kpi.Set(context.TODO(), _nsKey(config, pairIn.Host), string(val), nil)
			if err != nil {
				klog.Error(err)
			} else {
				klog.Infof("Updating DNS %s -> %s", pairIn.Host, val)
			}
		}
	}

	for _, pairOut := range out {
		// DELETE etcd
		kpi.Delete(context.TODO(), _nsKey(config, pairOut.Host), nil)
	}
}

// pushLease push self to etcd lease server
func pushLease(config *common.OverallConfig) {
	peer := models.PeerLease{
		Cidr:      config.Overlay.Cidr,
		Timestamp: time.Now().Unix(),
		TTL:       leaseTTL,
	}
	payload, err := json.Marshal(peer)
	if err != nil {
		klog.Error(err)
		return
	}

	kpi := clientv2.NewKeysAPI(etcd.NewClientOrDie())
	_, err = kpi.Set(context.TODO(), _leaseKey(config), string(payload), nil)
	if err != nil {
		klog.Error(err)
		return
	}
}

func pullLease(config *common.OverallConfig) ([]models.PeerLease, error) {
	res := make([]models.PeerLease, 0)
	kpi := clientv2.NewKeysAPI(etcd.NewClientOrDie())
	resp, err := kpi.Get(context.TODO(), config.Overlay.LeaseEtcdPrefix, nil)
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

func doLeaseSync(config *common.OverallConfig) {
	pushLease(config)
	peers, err := pullLease(config)
	if err != nil {
		klog.Errorf("Unable to pull lease %v", err)
	}
	for _, peer := range peers {
		if peer.Cidr == config.Overlay.Cidr {
			// Add route to bridge
			network.RouteToIface(peer.Cidr, config.BridgeName(), true)
		} else {
			if peer.Timestamp+peer.TTL > time.Now().Unix() {
				network.RouteToIface(peer.Cidr, flannel.FlannelIface, false)
			} else {
				// Delete expired routes
				network.DeleteRoute(peer.Cidr, flannel.FlannelIface)
			}
		}
	}
}

func doDomainSync(config *common.OverallConfig) {

	nodeManager := services.GetNodeManager(config.Backend)
	domains, err := nodeManager.List()
	if err != nil {
		klog.Error("Failed to list local domains", err)
		return
	}

	payload, err := json.Marshal(domains)
	if err != nil {
		klog.Error("Failed to marshal local domains", err)
	}
	kpi := clientv2.NewKeysAPI(etcd.NewClientOrDie())
	_, err = kpi.Set(context.TODO(), _domainKey(config), string(payload), nil)
	if err != nil {
		klog.Error("Failed to sync domains", err)
	}
}

func CommerceSync() {
	config := configurations.Config()
	go doLeaseSync(config)
	go doSync(config)
	go doDomainSync(config)
}
