package mesh

import (
	"context"
	"fmt"
	"mediumkube/pkg/common/flannel"
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/utils"
	"os"
	"os/exec"
	"strings"
	"time"

	etcd "mediumkube/pkg/etcd"

	clientv2 "go.etcd.io/etcd/client/v2"
	"k8s.io/klog/v2"
)

const (
	etcdExecutable = "etcd"
)

// start etcd service
func StartEtcd() *os.Process {
	etcdPort := configurations.Config().Overlay.EtcdPort
	master := configurations.Config().Overlay.Master
	cmd := exec.Command(
		"etcd",
		fmt.Sprintf("--listen-client-urls=%s", utils.EtcdEp(master, etcdPort)),
		fmt.Sprintf("--advertise-client-urls=%s", utils.EtcdEp(master, etcdPort)),
		"--enable-v2=true",
	)

	go utils.ExecWithStdio(cmd)
	time.Sleep(1 * time.Second)

	return cmd.Process
}

func initDnsDir() {
	klog.Info("Initialization etcd directory for dns")
	overlayConfig := configurations.Config().Overlay
	cli := etcd.NewClientOrDie()
	kpi := clientv2.NewKeysAPI(cli)
	kpi.Set(context.TODO(), overlayConfig.DNSEtcdPrefix, "", &clientv2.SetOptions{Dir: true})
}

// configFlannel render flannel configuration fron overall configurations
// and push to etcd
func configFlannel() {
	klog.Info("Initializing configurations for flannel")
	overlayConfig := configurations.Config().Overlay
	cli := etcd.NewClientOrDie()

	k := strings.Join([]string{overlayConfig.Flannel.EtcdPrefix, "config"}, "/")
	v := flannel.NewConfig(configurations.Config()).ToStr()
	kpi := clientv2.NewKeysAPI(cli)

	_, err := kpi.Set(context.TODO(), k, v, &clientv2.SetOptions{})
	klog.Info(k, v)
	if err != nil {
		klog.Error(err)
	}

}
