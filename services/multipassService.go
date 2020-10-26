package services

import (
	"fmt"
	"log"
	"mediumkube/utils"
	"os/exec"
)

type MultipassService struct{}

// Deploy deploy a kubernetes cluster
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
	cmd := "multipass launch -v -n node%v --cloud-init cloud-init.yaml -c %v -m %v -d %v %v"
	for i := 0; i < nodeNum; i++ {
		log.Printf("Deploying %v of %v nodes\n\r", i+1, nodeNum)
		thisCmd := fmt.Sprintf(cmd, i+1, cpu, mem, disk, img)
		execCmd := exec.Command(thisCmd)
		_, err := execCmd.Output()
		utils.CheckErr(err)
		log.Print("OK!")
	}
}
