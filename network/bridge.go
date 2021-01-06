package network

import (
	"fmt"
	"log"
	"mediumkube/common"
	"mediumkube/utils"

	"github.com/vishvananda/netlink"
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

	log.Println("Adding bridge: ", bridge.Name)
	err = netlink.LinkAdd(br)
	if err != nil {
		return err
	}

	lnk, err := netlink.LinkByName(bridge.Name)
	if err != nil {
		return err
	}

	log.Println("Setting up")
	err = netlink.LinkSetUp(lnk)
	if err != nil {
		return err
	}

	log.Println("Assigning address: ", addr)
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

// ShowBridge show details of bridge in stdout
func ShowBridge(bridge common.Bridge) {
	var lnk netlink.Link
	var err error

	if lnk, err = netlink.LinkByName(bridge.Name); err != nil {
		panic(err)
	}
	addrs, err := netlink.AddrList(lnk, ipv4)
	utils.CheckErr(err)
	fmt.Println("Name: ", lnk.Attrs().Name)
	fmt.Println("Addr: ", addrs)

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
