package network

import (
	"bufio"
	"context"
	"mediumkube/pkg/common"
	"mediumkube/pkg/utils"
	"os"
	"strings"

	clientv2 "go.etcd.io/etcd/client/v2"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

// time, mac, ip, name, clientID
func parse(leaseEntry string) (string, string, string, string, string) {
	fields := strings.Split(leaseEntry, " ")
	if len(fields) == 5 {
		return fields[0], fields[1], fields[2], fields[3], fields[4]
	}
	return "*", "*", "*", "*", "*"
}

// ListNSPairs list pairs from a local lease file
func ListNSPairs(leaseFilePath string) []common.NSPair {
	res := make([]common.NSPair, 0)
	file, err := os.Open(leaseFilePath)
	utils.CheckErr(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		entry := scanner.Text()
		_, _, ip, node, _ := parse(entry)
		if node == "*" {
			continue
		}
		res = append(res, common.NSPair{Host: node, Address: ip})
	}
	return res
}

// ListETCDNSPairs
// List pairs on etcd that matchs given cidr
// To fetch all pairs, set cidr to empty string
func ListETCDNsPairs(client clientv2.Client, dnsPrefix string, cidr string) []common.NSPair {
	res := make([]common.NSPair, 0)
	kpi := clientv2.NewKeysAPI(client)
	resp, err := kpi.Get(context.TODO(), dnsPrefix, nil)
	if err != nil {
		klog.Error(err)
		return res
	}
	for _, node := range resp.Node.Nodes {
		v := []byte(node.Value)
		pair := common.NSPair{}
		yaml.Unmarshal(v, &pair)
		if cidr == "" || utils.CidrMatch(pair.Address, cidr) {
			res = append(res, pair)
		}
	}
	return res
}

func exist(pairs []common.NSPair, tgt common.NSPair) bool {
	for _, pair := range pairs {
		if pair.Host == tgt.Host && pair.Address == tgt.Address {
			return true
		}
	}
	return false
}

func SyncDNSLease(local []common.NSPair, remote []common.NSPair) ([]common.NSPair, []common.NSPair) {
	in := make([]common.NSPair, 0)
	out := make([]common.NSPair, 0)
	for _, l := range local {
		if !exist(remote, l) {
			in = append(in, l)
		}
	}
	for _, r := range remote {
		if !exist(local, r) {
			out = append(out, r)

		}
	}

	return in, out
}

// Resolve node name to its ip address
func Resolve(leaseFilePath string, host string) (string, bool) {
	for _, nsPair := range ListNSPairs(leaseFilePath) {
		if nsPair.Host == host {
			return nsPair.Address, true
		}
	}

	return "0.0.0.0", false
}

func ResolveOverlay(overlayConfig common.Overlay, host string) (string, bool) {
	return "0.0.0.0", false
}
