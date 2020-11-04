package services

import (
	"fmt"
	"log"
	"mediumkube/common"
	"mediumkube/configurations"
	"mediumkube/utils"
	"os"
	"path/filepath"
	"time"

	"os/exec"

	"github.com/google/uuid"
)

// MultipassService interact with multipass using commands
type MultipassService struct {
	OverallConfig common.OverallConfig
}

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
func (service MultipassService) Deploy(nodeNum int, cpu string, mem string, disk string, img string, cloudInit string, mounts map[string]string) {
	for i := 0; i < nodeNum; i++ {
		log.Printf("Deploying %v of %v nodes", i+1, nodeNum)
		nodeName := fmt.Sprintf("node%v", i+1)

		go func() {
			for k, v := range mounts {
				os.MkdirAll(k, 0777)
				for {
					mountCmd := exec.Command("multipass", "mount", k, fmt.Sprintf("%v:%v", nodeName, v))
					_, err := mountCmd.Output()
					if err != nil {
						time.Sleep(5 * time.Second)
					}
				}
			}

		}()

		execCmd := exec.Command(
			"multipass",
			"launch",
			"-vvv",
			"-n", nodeName,
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

// ExecUnchecked Exec a command on a virtual machine without checking error
func (service MultipassService) ExecUnchecked(node string, command []string, sudo bool) (string, error) {

	return "", nil
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

func init() {
	InitMultipassService(MultipassService{
		OverallConfig: configurations.Config(),
	})
}
