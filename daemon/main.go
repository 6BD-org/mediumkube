package main

import (
	"flag"
	"fmt"
	"log"
	"mediumkube/common"
	"mediumkube/configurations"
	"mediumkube/network"
	"mediumkube/utils"
	"os"
	"time"

	"github.com/coreos/go-iptables/iptables"
	"github.com/vishvananda/netlink"
	"k8s.io/klog/v2"
)

const (
	ipv4  int    = netlink.FAMILY_V4
	chain string = "MEDIUMKUBE_FW"
	table string = "filter"
)

func _forwardRuleIn(bridge common.Bridge) []string {
	return []string{
		"-s", bridge.Inet,
		"-i", bridge.Name,
		"-j", "ACCEPT",
	}
}

func _forwardRuleOut(bridge common.Bridge) []string {
	return []string{
		"-d", bridge.Inet,
		"-o", bridge.Name,
		"-j", "ACCEPT",
	}
}

func _forwardRejectICMPUnreachableIn(bridge common.Bridge) []string {
	return []string{
		"-i", bridge.Name,
		"-j", "REJECT",
		"--reject-with", "icmp-port-unreachable",
	}
}

func _forwardRejectICMPUnreachableOut(bridge common.Bridge) []string {
	return []string{
		"-o", bridge.Name,
		"-j", "REJECT",
		"--reject-with", "icmp-port-unreachable",
	}
}

func _dhcpIn(bridge common.Bridge) []string {
	return []string{
		"-i", bridge.Name,
		"-p", "udp",
		"-m", "udp",
		"--dport", "67",
		"-j", "ACCEPT",
	}
}

func _dhcpOut(bridge common.Bridge) []string {
	return []string{
		"-o", bridge.Name,
		"-p", "udp",
		"-m", "udp",
		"--sport", "67",
		"-j", "ACCEPT",
	}
}

func addrsEq(addrs []netlink.Addr, bridgeAddr string) bool {
	if len(addrs) > 1 {
		// mediumkube only support single address
		return false
	}

	if len(addrs) == 1 {

	}

	return true
}

func refreshBridge() {
	// Step1: Find bridge defined in config
	bridgeName := configurations.Config().Bridge.Name
	bridgeAddr := configurations.Config().Bridge.Inet
	// hostNIC := configurations.Config().Bridge.Host

	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		log.Println(err)
	}
	addrs, err := netlink.AddrList(br, ipv4)
	if err != nil {
		log.Println(err)
	}

	log.Println("Get: ", addrs, "Expect: ", bridgeAddr)

}

func processExistence(bridge common.Bridge) {
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

func processAddr(bridge common.Bridge) {
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
		} else {
			klog.Info("Assigning address: ", newAddr)
			addErr := netlink.AddrAdd(lnk, newAddr)
			utils.WarnErr(addErr)
			return
		}
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

func createRuleIfNotExists(chain string, rules ...string) {
	iptable, err := iptables.New()
	if err != nil {
		klog.Error(err)
		return
	}

	exists, err := iptable.Exists(
		table,
		chain,
		rules...,
	)
	if err != nil {
		klog.Error(err)
		return
	}

	if !exists {
		klog.Info("Appending: ", rules)
		iptable.Append(table, chain, rules...)
	}
}

func processIptables(bridge common.Bridge) {
	createRuleIfNotExists("FORWARD", _forwardRuleOut(bridge)...)
	createRuleIfNotExists("FORWARD", _forwardRuleIn(bridge)...)
	createRuleIfNotExists("FORWARD", _forwardRejectICMPUnreachableIn(bridge)...)
	createRuleIfNotExists("FORWARD", _forwardRejectICMPUnreachableOut(bridge)...)
	createRuleIfNotExists("INPUT", _dhcpIn(bridge)...)
	createRuleIfNotExists("OUTPUT", _dhcpOut(bridge)...)
}

func main() {

	tmpFlagSet := flag.NewFlagSet("", flag.ExitOnError)
	configDir := tmpFlagSet.String("config", "./config.yaml", "Configuration file")
	tmpFlagSet.Parse(os.Args)
	configurations.InitConfig(*configDir)

	for {
		time.Sleep(5 * time.Second)
		bridge := configurations.Config().Bridge
		processExistence(bridge)
		processAddr(bridge)
		processIptables(bridge)
	}
}
