package services

import (
	"fmt"
	"log"
	"mediumkube/common"
	"mediumkube/configurations"
	"mediumkube/utils"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/libvirt/libvirt-go"
)

// LibvirtService is implementation of node manager
type LibvirtService struct {
	config *common.OverallConfig
	conn   *libvirt.Connect
}

// Deploy deploy a domain
func (service LibvirtService) Deploy(nodes []common.NodeConfig, cloudInit string, image string) {}

// Purge purge a domain
func (service LibvirtService) Purge(node string) {}

// Start start a domain
func (service LibvirtService) Start(node string) {}

// Stop stop a domain
func (service LibvirtService) Stop(node string) {}

// Exec a command in a domain and return output
func (service LibvirtService) Exec(node string, command []string, sudo bool) string {
	return ""
}

// Transfer a file to a domain
func (service LibvirtService) Transfer(src string, tgt string) {}

// AttachAndExec attach to std and execute
func (service LibvirtService) AttachAndExec(node string, command []string, sudo bool) {}

// ExecScript a script
func (service LibvirtService) ExecScript(node string, script string, sudo bool) {
	rndStr := uuid.New().String()
	targetDir := filepath.Join("/tmp", "mediumkube", "shell", rndStr)
	targetPath := filepath.Join(targetDir, utils.GetFileName(script))
	mkdirCmd := []string{"mkdir", "-p", targetDir}
	shCmd := []string{"sh", targetPath}
	rmCmd := []string{"rm", "-rf", targetDir}

	service.Exec(node, mkdirCmd, false)
	service.Transfer(script, fmt.Sprintf("%v:%v", node, targetPath))
	service.Exec(node, shCmd, sudo)

	log.Println("Shell execution finished! Cleaning up cache")
	service.Exec(node, rmCmd, false)
}

func init() {
	if configurations.Config().Backend == "libvirt" {
		conn, err := libvirt.NewConnect("qemu:///system")
		utils.CheckErr(err)
		InitLibvritService(
			LibvirtService{
				config: configurations.Config(),
				conn:   conn,
			},
		)
	}
}
