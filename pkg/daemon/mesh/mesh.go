package mesh

import (
	"os"

	"github.com/mitchellh/go-ps"
	"k8s.io/klog/v2"
)

var (
	meshProcesses map[*os.Process]bool
)

// HealthCheck checks Mesh daemon health
// returns flannelOk, etcdOk
func HealthCheck() (bool, bool) {
	processes, err := ps.Processes()
	if err != nil {
		klog.Error("Error occurred in mesh health check")
	}
	flannelOk := false
	etcdOk := false

	for _, p := range processes {
		if p.Executable() == flannelExecutableName {
			flannelOk = true
		}
		if p.Executable() == etcdExecutable {
			etcdOk = true
		}
	}
	return flannelOk, etcdOk

}

func StartMesh() {
	flannelOk, etcdOk := HealthCheck()

	if !etcdOk {
		etcdProc := StartEtcd()
		meshProcesses[etcdProc] = true
	}

	if !flannelOk {
		flannelProc := StartFlannel()
		meshProcesses[flannelProc] = true
		configFlannel()
	}

	if !etcdOk {
		// ETCD restarted, init dns dir
		initDnsDir()
	}

	StartDNSSync()
}

func StopMesh() {
	StopDNSSync()
	for k, _ := range meshProcesses {
		err := k.Kill()
		if err != nil {
			klog.Error(err)
		}
	}

}

func init() {
	meshProcesses = make(map[*os.Process]bool)
}
