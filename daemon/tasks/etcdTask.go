package tasks

import (
	"fmt"
	"mediumkube/configurations"
	"mediumkube/utils"
	"os"
	"os/exec"
	"time"
)

func StartEtcd() *os.Process {
	etcdPort := configurations.Config().Overlay.EtcdPort
	master := configurations.Config().Overlay.Master
	cmd := exec.Command(
		"etcd",
		fmt.Sprintf("--listen-client-urls=%s", utils.EtcdEp(master, etcdPort)),
		fmt.Sprintf("--advertise-client-urls=%s", utils.EtcdEp(master, etcdPort)),
		"--enable-v2=true",
	)

	go utils.ExecWithStdio(cmd)
	time.Sleep(1 * time.Second)

	return cmd.Process
}
