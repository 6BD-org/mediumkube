package services

import (
	"mediumkube/pkg/common"
	"mediumkube/pkg/models"
)

// NodeManager manages nodes
type NodeManager interface {
	Deploy(nodes []common.NodeConfig, image string, sink ...func([]byte) error)
	Purge(node string)
	Start(node string)
	Stop(node string)
	Exec(node string, command []string, sudo bool) string
	Transfer(src string, tgt string)
	TransferR(src string, tgt string)
	AttachAndExec(node string, command []string, sudo bool)
	Shell(node string)
	ExecScript(node string, script string, sudo bool)
	List() ([]models.Domain, error)
	Disconnect()
}
