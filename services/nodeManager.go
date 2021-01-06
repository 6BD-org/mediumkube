package services

import "mediumkube/common"

// NodeManager manages nodes
type NodeManager interface {
	Deploy(nodes []common.NodeConfig, cloudInit string, image string)
	Purge(node string)
	Start(node string)
	Stop(node string)
	Exec(node string, command []string, sudo bool) string
	Transfer(src string, tgt string)
	AttachAndExec(node string, command []string, sudo bool)
	ExecScript(node string, script string, sudo bool)
}
