package tasks

import (
	"fmt"
	"log"
	"mediumkube/common"
	"mediumkube/utils"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/mitchellh/go-ps"
	"github.com/vishvananda/netlink"
	"k8s.io/klog/v2"
)

func bridgeSubNet(bridge common.Bridge) string {
	ip := strings.Split(bridge.Inet, "/")[0]
	builder := strings.Builder{}
	ipSplitted := strings.Split(ip, ".")
	for i := 0; i < 3; i++ {
		builder.WriteString(ipSplitted[i])
		builder.WriteString(".")
	}
	res := builder.String()
	return res[:len(res)-1]
}

// prepare kills dnsmasq processes that are still running
func preapare() {
	processes, err := ps.Processes()
	utils.CheckErr(err)

	for _, p := range processes {
		if p.Executable() == "dnsmasq" {
			cmdLine := utils.GetLinuxProcCmdOrEmpty(p.Pid())
			if strings.Contains(cmdLine, "--domain=mediumkube") {
				klog.Info("Killing: ", cmdLine, "PID", p.Pid())
				osp, err := os.FindProcess(p.Pid())
				if err != nil {
					klog.Error("Error finding process: ", p.Pid())
				}
				errKill := osp.Kill()
				osp.Wait()
				if errKill != nil {
					klog.Error(errKill)
				}
			}

		}
	}
}

// StartDnsmasq for DNS and NAT
func StartDnsmasq(bridge common.Bridge, config common.OverallConfig) *os.Process {
	subnet := bridgeSubNet(bridge)
	for {
		_, err := netlink.LinkByName(bridge.Name)
		if err != nil {
			_, ok := err.(netlink.LinkNotFoundError)
			if ok {
				klog.Info("Waiting for bridge to be created")
			}
			log.Println(err)
		} else {
			break
		}
		time.Sleep(1 * time.Second)
	}

	leaseFile := path.Join(config.TmpDir, "dnsmasq.lease")

	cmd := exec.Command(
		"dnsmasq",
		"--keep-in-foreground",
		"--strict-order",
		"--bind-interfaces",
		"--pid-file",
		"--domain=mediumkube",
		"--local=/mediumkube/",
		"--except-interface=lo",
		"--interface", bridge.Name,
		fmt.Sprintf("--listen-address=%v", strings.Split(bridge.Inet, "/")[0]),
		"--dhcp-no-override",
		"--dhcp-authoritative",
		// NEVER change lease file.
		fmt.Sprintf("--dhcp-leasefile=%v", leaseFile),
		fmt.Sprintf("--dhcp-range=%v", fmt.Sprintf("%v.2,%v.254,infinite", subnet, subnet)),
	)
	preapare()
	klog.Info("Starting dnsmasq with: ", cmd)
	go utils.ExecWithStdio(cmd)
	time.Sleep(1 * time.Second)
	proc := cmd.Process
	return proc
}
