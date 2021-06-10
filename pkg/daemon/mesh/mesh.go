package mesh

import (
	"mediumkube/pkg/utils"
	"os"

	"github.com/mitchellh/go-ps"
	"k8s.io/klog/v2"
)

var (
	meshProcesses map[*os.Process]bool
)

// HealthCheck checks Mesh daemon health
// returns flannelOk, etcdOk
func HealthCheck() (*os.Process, *os.Process) {
	processes, err := ps.Processes()
	if err != nil {
		klog.Error("Error occurred in mesh health check")
	}
	var flannelP *os.Process = nil
	var etcdP *os.Process = nil
	for _, p := range processes {
		if utils.SameProcInThisContext(p.Executable(), flannelExecutableName) {
			flannelP, err = os.FindProcess(p.Pid())
		}
		if utils.SameProcInThisContext(p.Executable(), etcdExecutable) {
			etcdP, err = os.FindProcess(p.Pid())
		}
	}
	return flannelP, etcdP

}

// StartMesh is invoked repeatdly, so makesure everything inside this method
// is idempotent
func StartMesh() {
	flannelP, etcdP := HealthCheck()

	if etcdP == nil {
		klog.Info("Etcd process not detected, creating one")

		etcdProc := StartEtcd()
		meshProcesses[etcdProc] = true
	} else {
		meshProcesses[etcdP] = true
	}

	if flannelP == nil {
		klog.Info("Flannel process not detected, creating one")
		flannelProc := StartFlannel()
		meshProcesses[flannelProc] = true
		configFlannel()
	} else {
		meshProcesses[flannelP] = true
	}

	initDnsDir()

	CommerceSync()
}

func StopMesh() {
	for k, _ := range meshProcesses {
		klog.Infof("Killing: %v", k.Pid)
		err := k.Kill()
		if err != nil {
			klog.Error(err)
		}
	}

}

func init() {
	meshProcesses = make(map[*os.Process]bool)
}
