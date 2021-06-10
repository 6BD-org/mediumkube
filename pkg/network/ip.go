package network

import (
	"errors"
	"net"
)

// https://gist.github.com/kotakanbe/d3059af990252ba89a82
func Hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func CidrIPRange(cidr string) (string, string, error) {
	hosts, err := Hosts(cidr)
	if err != nil {
		return "", "", err
	}
	if len(hosts) < 2 {
		return "", "", errors.New("Insufficient hosts in cidr")
	}
	// ignore interface ip
	return hosts[1], hosts[len(hosts)-1], nil
}
