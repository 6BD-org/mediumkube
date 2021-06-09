package network

import (
	"fmt"
	"mediumkube/pkg/utils"
	"net"

	"github.com/vishvananda/netlink"
	"k8s.io/klog/v2"
)

// FindRoute
func FindRoute(cidrStr string, ifaceName string, self bool) (*netlink.Route, error) {
	var err error
ERROR:
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	_, cidr, err := net.ParseCIDR(cidrStr)
	if err != nil {
		goto ERROR
	}

	if lnk, err := netlink.LinkByName(ifaceName); err != nil {
		goto ERROR
	} else {
		if routes, err := netlink.RouteList(lnk, ipv4); err != nil {
			goto ERROR
		} else {
			for _, r := range routes {
				if utils.IpNetEqual(cidr, r.Dst) {
					return &r, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("Route not found for %v -> %v", cidrStr, ifaceName)

}

// RouteToIface create a route to interface if not exists.
func RouteToIface(cidrStr string, ifaceName string, self bool) (*netlink.Route, error) {
	var err error

ERROR:
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	flnk, err := netlink.LinkByName(ifaceName)
	if err != nil {
		goto ERROR
	}

	err = checkIfaceAddress(flnk)
	if err != nil {
		goto ERROR
	}

	route, err := constructRoute(cidrStr, ifaceName, self)
	if err != nil {
		goto ERROR
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
			goto ERROR
		}
	}

	return route, nil

}

func checkIfaceAddress(link netlink.Link) error {

	addrs, err := netlink.AddrList(link, ipv4)
	if err != nil {
		klog.Error(err)
		return err
	}
	if len(addrs) == 0 {
		return fmt.Errorf("Unable to fetch IP for flannel")
	}
	return nil
}

func constructRoute(cidrStr string, ifaceName string, self bool) (*netlink.Route, error) {
	flnk, err := netlink.LinkByName(ifaceName)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	_, cidr, err := net.ParseCIDR(cidrStr)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	var gateway net.IP

	if self {
		gateway = net.IPv4(0, 0, 0, 0)
	} else {
		gateway = cidr.IP
	}

	route := &netlink.Route{
		Dst:       cidr,
		LinkIndex: flnk.Attrs().Index,
		Gw:        gateway,
	}
	return route, nil
}
