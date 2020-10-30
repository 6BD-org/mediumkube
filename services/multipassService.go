package services

import (
	"fmt"
	"log"
	"mediumkube/utils"

	"os/exec"

	"github.com/google/uuid"
)

// MultipassService interact with multipass using commands
type MultipassService struct{}

// Deploy deploy a vm collection
// All the parameters are guaranteed to ne non-nil values from
// Upper layers
//
// The name of the nodes are node01, node02, ... nodeXX depending on how
// many nodes are deployed
//
// All the nodes are identical in terms of resources. This is
// for simplicity. If you are configuring a real cluster, think about QoS when
// configuring different nodes
//
// cpu is the number of cpu cores
// 	for example, 2
// mem is memory size in Gigabytes
// 	for example, 2G
// disk is disk space in Gigagytes
//	for example, 20G
// img can be either remote image name or local .img file. See multipass document for more details
// cloudInit is cloudInit file used by multipass
func (service MultipassService) Deploy(nodeNum int, cpu string, mem string, disk string, img string, cloudInit string) {
	for i := 0; i < nodeNum; i++ {
		log.Printf("Deploying %v of %v nodes", i+1, nodeNum)
		execCmd := exec.Command(
			"multipass",
			"launch",
			"-v",
			"-n", fmt.Sprintf("node%v", i+1),
			"--cloud-init", cloudInit,
			"-c", cpu,
			"-m", mem,
			"-d", disk,
			img,
		)
		_, err := utils.ExecWithStdio(execCmd)
		utils.CheckErr(err)
		log.Println("OK!")
	}
}

// KubeInit init k8s cluster on a node
func (service MultipassService) KubeInit(node string, command string) {
	log.Printf("Executing %v on node %v...", command, node)
	execCmd := exec.Command(
		"multipass",
		"exec",
		"-v",
		node,
		"--",
		command,
	)
	out, err := execCmd.Output()
	utils.CheckErr(err)
	log.Println(out)
}

func preExecPrecess(node string, command []string, sudo bool) *exec.Cmd {
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
	execCmd := preExecPrecess(node, command, sudo)
	out, err := utils.ExecWithStdio(execCmd)
	utils.CheckErr(err)
	log.Println(out)

	return out
}

// ExecUnchecked Exec a command on a virtual machine without checking error
func (service MultipassService) ExecUnchecked(node string, command []string, sudo bool) (string, error) {

	return "", nil
}

// AttachAndExec Execute a command on a virtual machine with stdio attached
func (service MultipassService) AttachAndExec(node string, command []string, sudo bool) {
	execCmd := preExecPrecess(node, command, sudo)
	utils.AttachAndExec(execCmd)
}

// ExecScript Execute a local script to a node
func (service MultipassService) ExecScript(node string, script string, sudo bool) {
	rndStr := uuid.New().String()
	targetDir := fmt.Sprintf("/tmp/mediumkube/shell/%v", rndStr)
	targetPath := fmt.Sprintf("%v/%v", targetDir, utils.GetFileName(script))

	mkdirCmd := []string{"mkdir", "-p", targetDir}
	sendCmd := exec.Command("multipass", "transfer", script, fmt.Sprintf("%v:%v", node, targetPath))
	shCmd := []string{"sh", targetPath}
	rmCmd := []string{"rm", "-rf", targetDir}

	service.Exec(node, mkdirCmd, false)

	utils.ExecWithStdio(sendCmd)
	service.Exec(node, shCmd, sudo)

	log.Println("Shell execution finished! Cleaning up cache")
	service.Exec(node, rmCmd, false)

}
