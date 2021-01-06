package services

import (
	"fmt"
	"log"
	"mediumkube/common"
	"mediumkube/configurations"
	"mediumkube/utils"
	"os"
	"os/exec"
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
func (service LibvirtService) Deploy(nodes []common.NodeConfig, cloudInit string, image string) {
	defer service.conn.Close()
	for _, n := range nodes {
		cloudImage := fmt.Sprintf("%v.iso", n.CDROM)
		os.Remove(cloudImage)
		cmdGenIso := exec.Command(
			"genisoimage", "-o", cloudImage,
			"-V", "cidata",
			"-r", "-J",
			"-V", "cidata",
			cloudInit,
		)
		utils.AttachAndExec(cmdGenIso)
		// Step1 Create disk
		err := os.RemoveAll(n.DiskImage)
		utils.CheckErr(err)
		cmd := exec.Command(
			"virt-install",
			"-n", n.Name,
			"--description", n.Name,
			"--os-type", "generic",
			// "--os-variant", "ubuntu",
			"--memory", fmt.Sprintf("%v", utils.Convert(n.MEM, utils.M)),
			"--vcpus", n.CPU,
			"--import",
			"--disk", fmt.Sprintf("path=%v", n.CDROM),
			"--disk", fmt.Sprintf("path=%v,device=cdrom", cloudImage),
			"--disk", fmt.Sprintf("path=%v,bus=%v,size=%v", n.DiskImage, "virtio", utils.Convert(n.DISK, utils.G)),
			"--network", fmt.Sprintf("bridge=%v,model=virtio", service.config.Bridge.Name),
			"--check", "path_in_use=off",
		)

		utils.AttachAndExec(cmd)

	}

}

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

func getVirtNetworkForNode(node common.NodeConfig, networks []libvirt.Network) libvirt.Network {
	for _, net := range networks {
		virtNetName, err := net.GetName()
		utils.CheckErr(err)
		if virtNetName == node.Network.Name {
			return net
		}
	}
	return libvirt.Network{}
}

// List domains
func (service LibvirtService) List() {
	defer service.conn.Close()
	config := configurations.Config()
	dms, err := service.conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	utils.CheckErr(err)
	networks, err := service.conn.ListAllNetworks(libvirt.CONNECT_LIST_NETWORKS_ACTIVE)
	utils.CheckErr(err)
	fmt.Println("Name \t IP \t Status \t Reason")
	for _, d := range dms {
		var node common.NodeConfig = common.NodeConfig{}
		var network libvirt.Network = libvirt.Network{}
		domainName, err := d.GetName()
		utils.CheckErr(err)

		for _, nd := range config.NodeConfig {
			if nd.Name == domainName {
				node = nd
				network = getVirtNetworkForNode(node, networks)
			}
		}
		netName, err := network.GetName()
		utils.CheckErr(err)
		domainState, r, err := d.GetState()
		utils.CheckErr(err)
		fmt.Printf("%v \t %v \t %v \t %v \n", node.Name, netName, domainState, r)
	}
}

func init() {
	log.Println("Initing socket connection")
	conn, err := libvirt.NewConnect("qemu:///system")
	utils.CheckErr(err)
	InitLibvritService(
		LibvirtService{
			config: configurations.Config(),
			conn:   conn,
		},
	)
}
