package common

import (
	"mediumkube/utils"
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
