package mesh

import (
	"mediumkube/pkg/common"
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/utils"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mitchellh/go-ps"
	"github.com/vishvananda/netlink"
	"k8s.io/klog/v2"
)

var (
	route *netlink.Route
)

const (
	ipv4                  int    = netlink.FAMILY_V4
	flannelExecutableName string = "flanneld"
)

func prepare(config *common.OverallConfig) {
	processes, err := ps.Processes()
	utils.CheckErr(err)

	for _, p := range processes {
		if p.Executable() == flannelExecutableName {
			cmd := utils.GetLinuxProcCmdOrEmpty(p.Pid())
			if strings.Contains(cmd, config.Overlay.Flannel.EtcdPrefix) {
				osp, err := os.FindProcess(p.Pid())
				if err != nil {
					klog.Error("Fail to find process for running flanneld")
					return
				}
				err = osp.Kill()
				if err != nil {
					klog.Error("Fail to kill existing flanneld")
					return
				}
				osp.Wait()
			}
		}
	}
}

func StartFlannel() *os.Process {
	etcdPort := configurations.Config().Overlay.EtcdPort
	master := configurations.Config().Overlay.Master
	cmd := exec.Command(
		"flanneld",
		"--etcd-endpoints", utils.EtcdEp(master, etcdPort),
		"--etcd-prefix", configurations.Config().Overlay.Flannel.EtcdPrefix,
		"--ip-masq",
	)
	prepare(configurations.Config())
	go utils.ExecWithStdio(cmd)

	time.Sleep(1 * time.Second)
	return cmd.Process
}

func ProcessRoute(bridge common.Bridge, flannel common.Flannel) {

}

func CleanUpRoute() {
	netlink.RouteDel(route)
}
