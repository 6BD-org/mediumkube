package mesh

import (
	"mediumkube/pkg/common"
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/utils"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/vishvananda/netlink"
	"k8s.io/klog/v2"
)

var (
	route *netlink.Route
)

func StartFlannel() *os.Process {
	etcdPort := configurations.Config().Overlay.EtcdPort
	master := configurations.Config().Overlay.Master
	cmd := exec.Command(
		"flanneld",
		"--etcd-endpoints", utils.EtcdEp(master, etcdPort),
		"--etcd-prefix", configurations.Config().Overlay.Flannel.EtcdPrefix,
		"--ip-masq",
	)

	go utils.ExecWithStdio(cmd)

	time.Sleep(1 * time.Second)
	return cmd.Process
}

func ProcessRoute(bridge common.Bridge, flannel common.Flannel) {

}

func addRouteToIface(cidrStr string, ifaceName string) {
	flnk, err := netlink.LinkByName(ifaceName)
	if err != nil {
		klog.Error(err)
		return
	}

	addrs, err := netlink.AddrList(flnk, ipv4)
	if err != nil {
		klog.Error(err)
		return
	}
	if len(addrs) == 0 {
		klog.Error("Unable to fetch IP for flannel")
	}

	_, cidr, err := net.ParseCIDR(cidrStr)
	if err != nil {
		klog.Error(err)
		return
	}

	route = &netlink.Route{
		Dst:       cidr,
		LinkIndex: flnk.Attrs().Index,
		Gw:        net.IPv4(0, 0, 0, 0),
	}

	routes, err := netlink.RouteList(flnk, ipv4)

	exists := false
	for _, r := range routes {
		if r.Dst == route.Dst && r.LinkIndex == route.LinkIndex {
			exists = true
		}
	}
	if !exists {
		// Create route if not exists
		err = netlink.RouteAdd(route)
		if err != nil {
			klog.Error(err)
			return
		}
	}

}

func CleanUpRoute() {
	netlink.RouteDel(route)
}
