package mesh

import (
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/utils"
	"os"
	"os/exec"
	"time"

	"github.com/vishvananda/netlink"
)

var (
	route *netlink.Route
)

const (
	ipv4                  int    = netlink.FAMILY_V4
	flannelExecutableName string = "mediumkube-flanneld"
)

func StartFlannel() *os.Process {
	etcdPort := configurations.Config().Overlay.EtcdPort
	master := configurations.Config().Overlay.Master
	cmd := exec.Command(
		flannelExecutableName,
		"--etcd-endpoints", utils.EtcdEp(master, etcdPort),
		"--etcd-prefix", configurations.Config().Overlay.Flannel.EtcdPrefix,
		"--ip-masq",
	)
	go utils.ExecWithStdio(cmd)

	time.Sleep(1 * time.Second)
	return cmd.Process
}
