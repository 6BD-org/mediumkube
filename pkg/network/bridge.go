package network

import (
	"mediumkube/pkg/common"
	"mediumkube/pkg/utils"

	"github.com/vishvananda/netlink"
	"k8s.io/klog/v2"
)

const (
	ipv4 int = 1
)

// CreateNetBridge from config
func CreateNetBridge(bridge common.Bridge) error {
	la := netlink.NewLinkAttrs()
	la.Alias = bridge.Alias
	la.Name = bridge.Name

	br := &netlink.Bridge{
		LinkAttrs: la,
	}
	addr, err := netlink.ParseAddr(bridge.Inet)
	if err != nil {
		return err
	}

	klog.Info("Adding bridge: ", bridge.Name)
	err = netlink.LinkAdd(br)
	if err != nil {
		return err
	}

	lnk, err := netlink.LinkByName(bridge.Name)
	if err != nil {
		return err
	}

	err = netlink.LinkSetUp(lnk)
	if err != nil {
		return err
	}

	klog.Info("Assigning address: ", addr)
	err = netlink.AddrAdd(lnk, addr)

	if err != nil {
		return err
	}

	return nil

}

// RemoveNetBridge Remove bridge defined in config
func RemoveNetBridge(bridge common.Bridge) error {

	var lnk netlink.Link
	var err error

	if lnk, err = netlink.LinkByName(bridge.Name); err != nil {
		return err
	}

	if err = netlink.LinkDel(lnk); err != nil {
		return err
	}

	return nil

}

/*
mpqemubr0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 10.135.114.1  netmask 255.255.255.0  broadcast 10.135.114.255
        inet6 fe80::439:e7ff:fe09:6d9a  prefixlen 64  scopeid 0x20<link>
        ether 06:39:e7:09:6d:9a  txqueuelen 1000  (Ethernet)
        RX packets 10662  bytes 869134 (869.1 KB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 25014  bytes 26252137 (26.2 MB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

*/

// ShowBridge show details of bridge in stdout
func ShowBridge(bridge common.Bridge) {
	var lnk netlink.Link
	var err error

	if lnk, err = netlink.LinkByName(bridge.Name); err != nil {
		panic(err)
	}
	_, err = netlink.AddrList(lnk, ipv4)
	utils.CheckErr(err)

}

// Up set bridge to up
func Up(bridge common.Bridge) error {
	lnk, err := netlink.LinkByName(bridge.Name)
	if err != nil {
		return err
	}

	err = netlink.LinkSetUp(lnk)
	if err != nil {
		return err
	}

	return nil
}

// Down set bridge to down
func Down(bridge common.Bridge) error {
	lnk, err := netlink.LinkByName(bridge.Name)
	if err != nil {
		return err
	}

	err = netlink.LinkSetDown(lnk)
	if err != nil {
		return err
	}

	return nil
}
