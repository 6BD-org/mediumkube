package tasks

import (
	"fmt"
	"log"
	"mediumkube/pkg/common"
	"mediumkube/pkg/network"
	"mediumkube/pkg/utils"

	"github.com/vishvananda/netlink"
	"k8s.io/klog/v2"
)

const (
	ipv4 int = netlink.FAMILY_V4
)

// ProcessExistence Create bridge if not exists
func ProcessExistence(bridge common.Bridge) {
	_, err := netlink.LinkByName(bridge.Name)
	if err != nil {
		_, ok := err.(netlink.LinkNotFoundError)
		if ok {
			network.CreateNetBridge(bridge)
			return
		}
		log.Println(err)
	}
}

// ProcessAddr Assign IP address to bridge
func ProcessAddr(bridge common.Bridge) {
	lnk, err := netlink.LinkByName(bridge.Name)
	if err != nil {
		klog.Error(err)
		return
	}

	addrs, err := netlink.AddrList(lnk, ipv4)
	if err != nil {
		klog.Error(err)
		return
	}
	newAddr, err := netlink.ParseAddr(bridge.Inet)
	if len(addrs) == 0 {
		// Create address

		if err != nil {
			klog.Error(err)
			return
		}
		klog.Info("Assigning address: ", newAddr)
		addErr := netlink.AddrAdd(lnk, newAddr)
		utils.WarnErr(addErr)
		return

	}

	addrsStrs := make([]string, len(addrs))
	for i, v := range addrs {
		size, _ := v.Mask.Size()
		addrsStrs[i] = fmt.Sprintf("%v/%v", v.IP.String(), size)
	}

	if !utils.Contains(addrsStrs, newAddr.String()) {
		klog.Info("Re-assigning address: ", addrsStrs, newAddr.String())
		for _, v := range addrs {
			err = netlink.AddrDel(lnk, &v)
			utils.WarnErr(err)
		}

		err = netlink.AddrAdd(lnk, newAddr)
		utils.WarnErr(err)
	}

}
