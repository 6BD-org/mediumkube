package utils

import (
	"crypto/rand"
	"net"
)

func GenerateMac() net.HardwareAddr {
	buf := make([]byte, 6)
	var mac net.HardwareAddr

	_, err := rand.Read(buf)
	if err != nil {
	}

	// Set the local bit
	buf[0] |= 2

	mac = append(mac, buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])

	return mac
}

func IpNetEqual(net1 *net.IPNet, net2 *net.IPNet) bool {
	return net1.String() == net2.String()
}
