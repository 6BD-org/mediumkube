package common

import (
	"fmt"
	"mediumkube/pkg/utils"
	"path"
)

// LoadFile load the content of of file into template
func (c OverallConfig) LoadFile(path string) string {
	return utils.ReadStr(path)
}

// PubKey vm public key into template
func (c OverallConfig) PubKey() string {
	return utils.ReadStr(c.PubKeyDir)
}

// PrivKey load vm private key into template
func (c OverallConfig) PrivKey() string {
	return utils.ReadStr(c.PrivKeyDir)
}

// HostPubKey load host public key into template
func (c OverallConfig) HostPubKey() string {
	return utils.ReadStr(c.HostPubKeyDir)
}

// HostPrivKey load host private key into template
func (c OverallConfig) HostPrivKey() string {
	return utils.ReadStr(c.HostPrivKeyDir)
}

// NodeCPU Get CPU of a node
func (c OverallConfig) NodeCPU(node string) string {
	for _, nc := range c.NodeConfig {
		if nc.Name == node {
			return nc.CPU
		}
	}
	panic("Node not found")
}

// NodeMemory Get Memory of a node
func (c OverallConfig) NodeMemory(node string) string {
	for _, nc := range c.NodeConfig {
		if nc.Name == node {
			return nc.MEM
		}
	}
	panic("Node not found")
}

// NodeDiskImage Get disk image path of a node
func (c OverallConfig) NodeDiskImage(node string) string {
	return path.Join(c.TmpDir, fmt.Sprintf("%v-os.img", node))
}

// BridgeName get mediumkube bridge name
func (c OverallConfig) BridgeName() string {
	return c.Bridge.Name
}
