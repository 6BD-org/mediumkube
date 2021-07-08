package etcd

import (
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/utils"

	clientv2 "go.etcd.io/etcd/client/v2"
)

var (
	client *clientv2.Client = nil
)

// NewClientOrDie get client. Created for first time
func NewClientOrDie() clientv2.Client {
	if client == nil {
		overlayConfig := configurations.Config().Overlay
		cli, err := clientv2.New(
			clientv2.Config{
				Endpoints: []string{
					utils.EtcdEp(overlayConfig.Master, overlayConfig.EtcdPort),
				},
			},
		)
		utils.CheckErr(err)
		client = &cli
	}

	return *client
}
