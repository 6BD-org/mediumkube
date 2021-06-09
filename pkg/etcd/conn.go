package etcd

import (
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/utils"

	clientv2 "go.etcd.io/etcd/client/v2"
)

func NewClientOrDie() clientv2.Client {
	overlayConfig := configurations.Config().Overlay
	cli, err := clientv2.New(
		clientv2.Config{
			Endpoints: []string{
				utils.EtcdEp(overlayConfig.Master, overlayConfig.EtcdPort),
			},
		},
	)
	utils.CheckErr(err)

	return cli
}
