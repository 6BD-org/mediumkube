package mesh

import (
	"context"
	"encoding/json"
	"mediumkube/pkg/configurations"
	etcd "mediumkube/pkg/etcd"
	"mediumkube/pkg/network"
	"sync"
	"time"

	clientv2 "go.etcd.io/etcd/client/v2"
	"k8s.io/klog/v2"
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

func DSNSyncD() {
	etcdClient = etcd.NewClientOrDie()
	for on {
		doSync()
		time.Sleep(3 * time.Second)
	}
}

func StartDNSSync() {
	mux.Lock()
	defer mux.Unlock()
	on = true
	go DSNSyncD()
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
