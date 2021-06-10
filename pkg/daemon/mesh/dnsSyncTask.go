package mesh

import (
	"context"
	"encoding/json"
	"mediumkube/pkg/common"
	"mediumkube/pkg/common/flannel"
	"mediumkube/pkg/configurations"
	etcd "mediumkube/pkg/etcd"
	"mediumkube/pkg/network"
	"strings"
	"sync"
	"time"

	clientv2 "go.etcd.io/etcd/client/v2"
	"k8s.io/klog/v2"
)

const (
	leaseTTL = 60 // 6 seconds
)

var (
	etcdClient clientv2.Client
	on         bool
	mux        sync.Mutex
)

func doSync() {
	config := configurations.Config()
	locals := network.ListNSPairs(config.LeaseFile())
	remotes := network.ListETCDNsPairs(etcdClient, config.Overlay.DNSEtcdPrefix, config.Overlay.Cidr)
	in, out := network.SyncDNSLease(locals, remotes)

	prefix := config.Overlay.DNSEtcdPrefix
	kpi := clientv2.NewKeysAPI(etcdClient)

	for _, pairIn := range in {
		// PUT etcd
		val, marshalErr := json.Marshal(&pairIn)
		prefixedKey := prefix + "/" + pairIn.Host
		if marshalErr != nil {
			klog.Error(marshalErr)
		} else {
			_, err := kpi.Set(context.TODO(), prefixedKey, string(val), nil)
			if err != nil {
				klog.Error(err)
			} else {
				klog.Infof("Updating DNS %s -> %s", pairIn.Host, val)
			}
		}
	}

	for _, pairOut := range out {
		// DELETE etcd
		kpi.Delete(context.TODO(), pairOut.Host, nil)
	}
}

func dnsSyncD() {
	for on {
		doSync()
		time.Sleep(3 * time.Second)
	}
}

// pushLease push self to etcd lease server
func pushLease(config *common.OverallConfig) {
	peer := PeerLease{
		Cidr:      config.Overlay.Cidr,
		Timestamp: time.Now().Unix(),
		TTL:       leaseTTL,
	}
	payload, err := json.Marshal(peer)
	if err != nil {
		klog.Error(err)
		return
	}

	kpi := clientv2.NewKeysAPI(etcdClient)
	key := config.Overlay.LeaseEtcdPrefix + "/" + strings.Replace(config.Overlay.Cidr, "/", "-", -1)
	_, err = kpi.Set(context.TODO(), key, string(payload), nil)
	if err != nil {
		klog.Error(err)
		return
	}
}

func pullLease(config *common.OverallConfig) ([]PeerLease, error) {
	res := make([]PeerLease, 0)
	kpi := clientv2.NewKeysAPI(etcdClient)
	resp, err := kpi.Get(context.TODO(), config.Overlay.LeaseEtcdPrefix, nil)
	if err != nil {
		klog.Error(err)
		return []PeerLease{}, err
	}

	for _, node := range resp.Node.Nodes {
		payload := PeerLease{}
		if len(node.Value) == 0 {
			continue
		}
		err = json.Unmarshal([]byte(node.Value), &payload)
		if err != nil {
			klog.Error(err)
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

func leaseSyncD() {
	config := configurations.Config()
	for on {
		doLeaseSync(config)
		time.Sleep(3 * time.Second)
	}
}

func StartDNSSync() {
	mux.Lock()
	defer mux.Unlock()
	on = true
	if etcdClient == nil {
		etcdClient = etcd.NewClientOrDie()
	}

	go dnsSyncD()
	go leaseSyncD()
}

func StopDNSSync() {
	klog.Info("Stopping dns sync")
	mux.Lock()
	defer mux.Unlock()
	on = false
}

func init() {
	on = false
	mux = sync.Mutex{}
}
