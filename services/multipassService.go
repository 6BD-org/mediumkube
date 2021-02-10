package services

import (
	"fmt"
	"log"
	"mediumkube/common"
	"mediumkube/configurations"
	"mediumkube/utils"
	"path/filepath"

	"os/exec"

	"github.com/google/uuid"
)

// MultipassService interact with multipass using commands
type MultipassService struct {
	OverallConfig *common.OverallConfig
}

// Deploy deploy a vm collection
// All the parameters are guaranteed to ne non-nil values from
// Upper layers
//
// All the nodes are identical in terms of resources. This is
// for simplicity. If you are configuring a real cluster, think about QoS when
// configuring different nodes
//
// image can be either remote image name or local .img file. See multipass document for more details
// cloudInit is cloudInit file used by multipass
func (service MultipassService) Deploy(nodes []common.NodeConfig, cloudInit string, image string) {
	for _, node := range nodes {
		log.Printf("Deploying node: %v", node.Name)

		execCmd := exec.Command(
			"multipass",
			"launch",
			"-v",
			"-v",
			"-v",
			"-n", node.Name,
			"--cloud-init", cloudInit,
			"-c", node.CPU,
			"-m", node.MEM,
			"-d", node.DISK,
			image,
		)
		_, err := utils.ExecWithStdio(execCmd)
		utils.CheckErr(err)
		log.Println("Deploy successfully")
	}
}

func preExecProcess(node string, command []string, sudo bool) *exec.Cmd {
	log.Printf("Executing %v on node %v...", command, node)
	if sudo {
		// sudo inside the virtual machine, so prepend
		command = append([]string{"sudo"}, command...)
	}
	command = append([]string{"exec", "-v", node, "--"}, command...)
	execCmd := exec.Command(
		"multipass", command...,
	)

	return execCmd
}

// Exec Execute a command on a virtual machine
func (service MultipassService) Exec(node string, command []string, sudo bool) string {
	execCmd := preExecProcess(node, command, sudo)
	out, err := utils.ExecWithStdio(execCmd)
	utils.CheckErr(err)
	log.Println(out)

	return out
}

// Transfer transfer file between vm and host machine
func (service MultipassService) Transfer(src string, tgt string) {
	transferCmd := exec.Command("multipass", "transfer", src, tgt)
	utils.ExecWithStdio(transferCmd)
}

// AttachAndExec Execute a command on a virtual machine with stdio attached
func (service MultipassService) AttachAndExec(node string, command []string, sudo bool) {
	execCmd := preExecProcess(node, command, sudo)
	utils.AttachAndExec(execCmd)
}

// ExecScript Execute a local script to a node
func (service MultipassService) ExecScript(node string, script string, sudo bool) {
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

// Purge delegated to multipass
func (service MultipassService) Purge(node string) {}

// Start delegated to
func (service MultipassService) Start(node string) {}

// Stop delegated to multipass
func (service MultipassService) Stop(node string) {}

// List delegated to multipass
func (service MultipassService) List() {}

// Shell delegated to multipass
func (service MultipassService) Shell(node string) {}

func init() {
	InitMultipassService(MultipassService{
		OverallConfig: configurations.Config(),
	})
}
