package tasks

import (
	"mediumkube/configurations"
	"mediumkube/utils"
	"os"
	"os/exec"
	"time"
)

func StartFlannel() *os.Process {
	etcdPort := configurations.Config().Overlay.EtcdPort
	master := configurations.Config().Overlay.Master
	cmd := exec.Command(
		"flanneld",
		"--etcd-endpoints", utils.EtcdEp(master, etcdPort),
		"--etcd-prefix", configurations.Config().Overlay.Flannel.EtcdPrefix,
	)

	go utils.ExecWithStdio(cmd)

	time.Sleep(1 * time.Second)
	return cmd.Process
}
