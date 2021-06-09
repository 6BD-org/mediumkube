package network

import (
	"bufio"
	"mediumkube/pkg/common"
	"mediumkube/pkg/utils"
	"os"
	"strings"
)

// time, mac, ip, name, clientID
func parse(leaseEntry string) (string, string, string, string, string) {
	fields := strings.Split(leaseEntry, " ")
	return fields[0], fields[1], fields[2], fields[3], fields[4]
}

// Resolve node name to its ip address
func Resolve(leaseFilePath string, host string) (string, bool) {
	file, err := os.Open(leaseFilePath)
	utils.CheckErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		entry := scanner.Text()
		_, _, ip, node, _ := parse(entry)
		if node == host {
			return ip, true
		}
	}

	return "0.0.0.0", false
}

func ResolveOverlay(overlayConfig common.Overlay, host string) (string, bool) {
	return "0.0.0.0", false
}
