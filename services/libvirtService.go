package services

import (
	"fmt"
	"log"
	"mediumkube/common"
	"mediumkube/configurations"
	"mediumkube/utils"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/libvirt/libvirt-go"
)

const (
	// 	networkTemplate string = `
	// <network>
	// 	<name>%v</name>
	// 	<forward/>
	// 	<bridge name='%v'/>
	// 	<ip address='%v.1' netmask='255.255.255.0'>
	// 		<dhcp>
	// 			<range start='%v.2' end='%v.254'></range>
	// 		</dhcp>
	// 	</ip>
	// </network>
	// `

	cloudInitMetaTemplate string = `
instance-id: id-%v
local-hostname: %v`
)

var (
	fileToClean = make([]string, 0)
)

func cleanUp() {
	for _, file := range fileToClean {
		os.RemoveAll(file)
	}
}

func bridgeIP(bridge common.Bridge) string {
	return strings.Split(bridge.Inet, "/")[0]
}

// func networkXML(name string, bridge common.Bridge) string {
// 	return fmt.Sprintf(
// 		networkTemplate,
// 		name,
// 		bridge.Name,
// 		bridgeSubNet(bridge),
// 		bridgeSubNet(bridge),
// 		bridgeSubNet(bridge),
// 	)
// }

func meta(nodeName string) string {
	return fmt.Sprintf(cloudInitMetaTemplate, nodeName, nodeName)
}

// func (service LibvirtService) createNetwork(name string, bridge common.Bridge) error {
// 	netDestroyCmd := exec.Command(
// 		"virsh",
// 		"net-destroy",
// 		name,
// 	)
// 	utils.ExecWithStdio(netDestroyCmd)

// 	xml := networkXML(name, bridge)
// 	log.Println(xml)
// 	_, err := service.conn.NetworkCreateXML(xml)
// 	return err
// }

func (service LibvirtService) createCloudInitCD(cloudInit string, nodeName string) string {

	userData := path.Join(service.config.TmpDir, "user-data")
	metaData := path.Join(service.config.TmpDir, "meta-data")

	utils.Copy(cloudInit, userData)
	utils.WriteStr(metaData, meta(nodeName), os.FileMode(0666))

	cloudImage := path.Join(service.config.TmpDir, fmt.Sprintf("%v-cloudinit.iso", nodeName))
	os.Remove(cloudImage)
	cmdGenIso := exec.Command(
		"genisoimage", "-o", cloudImage,
		"-V", "cidata",
		"-r", "-J",
		"-V", "cidata",
		userData, metaData,
	)
	utils.AttachAndExec(cmdGenIso)
	fileToClean = append(fileToClean, userData)
	fileToClean = append(fileToClean, metaData)
	return cloudImage
}

func copyAndResizeMedia(src string, tgt string, size string) {

	utils.Copy(src, tgt)

	cmd := exec.Command(
		"qemu-img", "resize",
		tgt, fmt.Sprintf("+%v", size),
	)
	_, err := utils.ExecWithStdio(cmd)
	utils.CheckErr(err)
}

// createDomain Create a domain, overwriting disk image
func (service LibvirtService) createDomain(name string, cpu string, memory string, disk string, net string, image string, cloudInitImg string) {
	// Step1 Create disk
	diskImage := path.Join(service.config.TmpDir, fmt.Sprintf("%v-disk.img", name))
	err := os.RemoveAll(diskImage)
	utils.CheckErr(err)
	cmd := exec.Command(
		"virt-install",
		"-n", name,
		"--os-type", "generic",
		// "--os-variant", "ubuntu",
		"--memory", fmt.Sprintf("%v", utils.Convert(memory, utils.M)),
		"--vcpus", cpu,
		"--import",
		"--disk", fmt.Sprintf("path=%v", image),
		"--disk", fmt.Sprintf("path=%v,device=cdrom", cloudInitImg),
		// "--disk", fmt.Sprintf("path=%v,bus=%v,size=%v", diskImage, "virtio", ,
		"--network", fmt.Sprintf("bridge=%v", net),
		"--check", "path_in_use=off",
		"--nographics",
	)
	utils.AttachAndExec(cmd)
}

// LibvirtService is implementation of node manager
type LibvirtService struct {
	config *common.OverallConfig
	conn   *libvirt.Connect
}

// Deploy deploy a domain
func (service LibvirtService) Deploy(nodes []common.NodeConfig, cloudInit string, image string) {
	defer service.conn.Close()
	defer cleanUp()
	for _, n := range nodes {
		// Step 0: Cloud init iso
		cloudImage := service.createCloudInitCD(cloudInit, n.Name)

		// Step2 Copy image
		srcImg := service.config.Image
		tgtImg := path.Join(service.config.TmpDir, fmt.Sprintf("%v-os.img", n.Name))
		log.Println("Copying image file from", srcImg, "to", tgtImg)

		copyAndResizeMedia(srcImg, tgtImg, n.DISK)

		// Step3 Create domain
		log.Println("Launching domain...")
		service.createDomain(
			n.Name,
			n.CPU,
			n.MEM,
			n.DISK,
			service.config.Bridge.Name,
			tgtImg,
			cloudImage,
		)
	}

}

// Purge purge a domain
func (service LibvirtService) Purge(node string) {
	// Step1 destory
	cmdDestory := exec.Command(
		"virsh",
		"destroy",
		node,
	)

	_, err := utils.ExecWithStdio(cmdDestory)
	utils.CheckErr(err)
	// Step2 undefine
	cmdUndefine := exec.Command(
		"virsh",
		"undefine",
		node,
		"--remove-all-storage",
	)
	_, err = utils.ExecWithStdio(cmdUndefine)
	utils.CheckErr(err)

}

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

// List domains
func (service LibvirtService) List() {
	defer service.conn.Close()
	dms, err := service.conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	utils.CheckErr(err)
	fmt.Println("Name \t IP \t Status \t Reason")
	for _, d := range dms {
		var node common.NodeConfig = common.NodeConfig{}

		domainState, r, _ := d.GetState()
		fmt.Printf("%v \t %v \t %v \n", node.Name, domainState, r)
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
