package services

import (
	"bytes"
	"fmt"
	"log"
	"mediumkube/pkg/common"
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/mediumssh"
	"mediumkube/pkg/models"
	"mediumkube/pkg/utils"
	"mediumkube/pkg/utils/virtutils"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/libvirt/libvirt-go"
	"k8s.io/klog/v2"
)

const (
	cloudInitMetaTemplate string = `
instance-id: id-%v
local-hostname: %v`

	sshUser   string  = "ubuntu"
	sshPort   int     = 22
	STDOUT_FD uintptr = 0
)

var (
	fileToClean = make([]string, 0)
	stateMap    = make(map[libvirt.DomainState]string)
)

func cleanUp() {
	for _, file := range fileToClean {
		os.RemoveAll(file)
	}
}

func bridgeIP(bridge common.Bridge) string {
	return strings.Split(bridge.Inet, "/")[0]
}

func formatSSHAddr(addr string) string {
	return fmt.Sprintf("%v:%v", addr, sshPort)
}

func (service LibvirtService) connectToNode(node string) (*mediumssh.SSHClient, error) {
	addr, ok := GetDNSService(service.config).Resolve(node)
	if !ok {
		klog.Error("Unable to resolve node: ", node)
		return nil, fmt.Errorf("Unable to resolve node: %v", node)
	}
	addr = formatSSHAddr(addr)
	sshClient := mediumssh.SSHLogin(sshUser, addr, service.config.HostPrivKeyDir)
	return sshClient, nil
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

// CreateDomain Create a domain, overwriting disk image
func (service LibvirtService) CreateDomain(name string, cpu string, memory string, disk string, net string, image string, cloudInitImg string) {
	param := common.NewDomainCreationParam(
		name, cpu, memory, image, cloudInitImg, net,
	)
	xml, err := virtutils.GetDeploymentConfig(
		param,
	)
	if err != nil {
		klog.Error(err)
		return
	}

	domain, err := service.conn.DomainDefineXMLFlags(xml, libvirt.DOMAIN_DEFINE_VALIDATE)
	if err != nil {
		klog.Error(err)
		return
	}

	err = domain.Create()
	if err != nil {
		klog.Error(err)
		return
	}

	bufStr := bytes.NewBufferString("")
	var handleWatch int = -1

	mux := sync.Mutex{}
	f := os.NewFile(STDOUT_FD, "")

	steamOut := func(stream *libvirt.Stream, cbType libvirt.StreamEventType) {
		mux.Lock()
		defer mux.Unlock()

		if cbType&libvirt.STREAM_EVENT_READABLE > 0 {

			for {
				ioBuffer := make([]byte, 1024)
				n, err := stream.Recv(ioBuffer)
				if n <= 0 {
					break
				}
				if err != nil {
					break
				}
				bufStr.Write(ioBuffer[:n])
			}
			if bufStr.Len() > 0 {
				libvirt.EventUpdateHandle(handleWatch, libvirt.EVENT_HANDLE_WRITABLE)
			}
		}
	}

	eventHandle := func(watch int, file int, events libvirt.EventHandleType) {
		mux.Lock()
		defer mux.Unlock()
		if events&libvirt.EVENT_HANDLE_WRITABLE > 0 {
			f.Write(bufStr.Bytes())
			bufStr.Reset()
			libvirt.EventUpdateHandle(handleWatch, 0)
		}
	}

	stream, err := service.conn.NewStream(libvirt.STREAM_NONBLOCK)

	err = domain.OpenConsole("", stream, 0)
	if err != nil {
		klog.Error(err)
		return
	}

	handleWatch, err = libvirt.EventAddHandle(1, 0, eventHandle)
	if err != nil {
		klog.Error(err)
		return
	}
	if handleWatch < 0 {
		klog.Error("Unable to register event handler")
	}

	err = stream.EventAddCallback(libvirt.STREAM_EVENT_READABLE, steamOut)
	if err != nil {
		klog.Error(err)
		return
	}

	for {
		if err := libvirt.EventRunDefaultImpl(); err != nil {
			klog.Error(err)
		}
	}
}

// LibvirtService is implementation of node manager
type LibvirtService struct {
	config *common.OverallConfig
	conn   *libvirt.Connect
}

// Deploy deploy a domain
// In libvirt backend, remote images are nolonger supported.
func (service LibvirtService) Deploy(nodes []common.NodeConfig, cloudInit string, image string) {
	defer service.conn.Close()
	defer cleanUp()
	for _, n := range nodes {
		// Step 0: Cloud init iso
		cloudImage := service.createCloudInitCD(cloudInit, n.Name)

		// Step2 Copy image
		srcImg := service.config.Image
		tgtImg := service.config.NodeDiskImage(n.Name)
		log.Println("Copying image file from", srcImg, "to", tgtImg)

		copyAndResizeMedia(srcImg, tgtImg, n.DISK)

		// Step3 Create domain
		log.Println("Launching domain...")
		service.CreateDomain(
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
// If the domain is running, this command will stop it
// then delete the domain along with storages attached to it
func (service LibvirtService) Purge(node string) {

	domain, err := service.conn.LookupDomainByName(node)
	if err != nil {
		klog.Error(err)
		return
	}

	// Destroy the domain if it is running
	if state, _, err := domain.GetState(); err == nil && state == libvirt.DOMAIN_RUNNING {
		klog.Info("Stopping node", node)
		err := domain.Destroy()
		if err != nil {
			klog.Error(err)
		}
	}

	err = domain.Undefine()
	if err != nil {
		klog.Error(err)
	}

	err = os.Remove(service.config.NodeDiskImage(node))
	if err != nil {
		klog.Error(err)
	}
}

// Start start a domain
func (service LibvirtService) Start(node string) {
	domain, err := service.conn.LookupDomainByName(node)
	if err != nil {
		klog.Error(err)
		return
	}
	err = domain.Create()
	if err != nil {
		klog.Error(err)
	}
}

// Stop stop a domain gracefully
func (service LibvirtService) Stop(node string) {
	domain, err := service.conn.LookupDomainByName(node)
	if err != nil {
		klog.Error(err)
		return
	}
	err = domain.DestroyFlags(libvirt.DOMAIN_DESTROY_GRACEFUL)
	if err != nil {
		klog.Error(err)
	}
}

// Exec a command in a domain and return output
func (service LibvirtService) Exec(node string, command []string, sudo bool) string {

	sshClient, err := service.connectToNode(node)
	utils.CheckErr(err)
	out := sshClient.Execute(command, sudo)
	fmt.Println(out)
	return out
}

// TransferR is recursive version of Transfer
func (service LibvirtService) TransferR(hostAndSrc string, hostAndTgt string) {
	if strings.Contains(hostAndTgt, ":") {
		hostTgt := strings.Split(hostAndTgt, ":")
		if len(hostTgt) < 2 {
			klog.Error("Invalid argument")
			return
		}
		host, tgt := hostTgt[0], hostTgt[1]

		src := hostAndSrc
		sshClient, err := service.connectToNode(host)
		utils.CheckErr(err)

		sshClient.TransferR(src, tgt)
		return
	}
	if strings.Contains(hostAndSrc, ":") {
		hostSrc := strings.Split(hostAndSrc, ":")
		if len(hostSrc) < 2 {
			klog.Error("Invalid argument")
			return
		}
		host, src := hostSrc[0], hostSrc[1]
		tgt := hostAndTgt

		sshClient, err := service.connectToNode(host)
		utils.CheckErr(err)
		sshClient.ReceiveR(src, tgt)
		return
	}
	klog.Error("Invalid argument")

}

// Transfer a file between vm and local machine
func (service LibvirtService) Transfer(hostAndSrc string, hostAndTgt string) {
	if strings.Contains(hostAndTgt, ":") {
		hostTgt := strings.Split(hostAndTgt, ":")
		if len(hostTgt) < 2 {
			klog.Error("Invalid argument")
			return
		}
		host, tgt := hostTgt[0], hostTgt[1]

		src := hostAndSrc
		sshClient, err := service.connectToNode(host)
		utils.CheckErr(err)

		sshClient.Transfer(src, tgt)
		return
	}

	if strings.Contains(hostAndSrc, ":") {
		hostSrc := strings.Split(hostAndSrc, ":")
		if len(hostSrc) < 2 {
			klog.Error("Invalid argument")
			return
		}
		host, src := hostSrc[0], hostSrc[1]
		tgt := hostAndTgt

		sshClient, err := service.connectToNode(host)
		utils.CheckErr(err)
		sshClient.Receive(src, tgt)
		return
	}

	klog.Error("Invalid argument")

}

// AttachAndExec attach to std and execute
func (service LibvirtService) AttachAndExec(node string, command []string, sudo bool) {
	sshClient, err := service.connectToNode(node)
	utils.CheckErr(err)
	sshClient.AttachAndExecute(command, sudo)
}

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

// Shell launch a ssh session to a domain
func (service LibvirtService) Shell(node string) {
	sshClient, err := service.connectToNode(node)
	utils.CheckErr(err)
	sshClient.Shell()
}

// List domains
func (service LibvirtService) List() ([]models.Domain, error) {
	defer service.conn.Close()
	dms, err := service.conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_PERSISTENT)

	utils.CheckErr(err)

	res := make([]models.Domain, 0)
	for _, d := range dms {
		name, _ := d.GetName()
		domainState, r, err := (&d).GetState()
		stateStr := ""
		if err != nil {
			klog.Error(err)
			return make([]models.Domain, 0), err
		} else {
			stateStr = stateMap[domainState]
		}
		addr, ok := GetDNSService(service.config).Resolve(name)
		if !ok {
			addr = "UNAVAILABLE"
		}

		res = append(res, models.Domain{
			Name:   name,
			IP:     addr,
			Status: stateStr,
			Reason: fmt.Sprint(r),
		})
	}
	return res, nil
}

func init() {
	log.Println("Initing socket connection")
	libvirt.EventRegisterDefaultImpl()
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		klog.Error("Fail to connect to libvirt: ", err)
	}
	InitLibvritService(
		LibvirtService{
			config: configurations.Config(),
			conn:   conn,
		},
	)

	stateMap[libvirt.DOMAIN_NOSTATE] = "NOSTATE"
	stateMap[libvirt.DOMAIN_RUNNING] = "RUNNING"
	stateMap[libvirt.DOMAIN_BLOCKED] = "BLOCKED"
	stateMap[libvirt.DOMAIN_PAUSED] = "PAUSED"
	stateMap[libvirt.DOMAIN_SHUTDOWN] = "SHUTDOWN"
	stateMap[libvirt.DOMAIN_CRASHED] = "CRASHED"
	stateMap[libvirt.DOMAIN_PMSUSPENDED] = "PMSUSPENDED"
	stateMap[libvirt.DOMAIN_SHUTOFF] = "SHUTOFF"
}
