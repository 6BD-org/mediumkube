package plugins

import (
	"fmt"
	"mediumkube/configurations"
	"mediumkube/mediumssh"
	"mediumkube/network"

	"k8s.io/klog/v2"
)

const (
	sshUser string = "ubuntu"
	sshPort int    = 22
)

type PostInitPlugin struct{}

func (plugin PostInitPlugin) Desc() {
	fmt.Println("Some jobs after kuberentes master is initialized")

}

func (plugin PostInitPlugin) Exec(args ...string) {
	if len(args) < 1 {
		klog.Error("Invalid Argument")
		return
	}
	host, ok := network.Resolve(configurations.Config().LeaseFile(), args[0])
	if !ok {
		klog.Error("Unknown host: ", args[0])
	}
	host = fmt.Sprintf("%v:%v", host, sshPort)
	sshClient := mediumssh.SSHLogin(
		sshUser, host, configurations.Config().HostPrivKeyDir,
	)

	sshClient.Execute(
		[]string{"mkdir", "-p", "$HOME/.kube"},
		false,
	)

	sshClient.Execute(
		[]string{"cp", "-i", "/etc/kubernetes/admin.conf", "$HOME/.kube/config"},
		true,
	)

	sshClient.Execute(
		[]string{"chown", "$(id -u):$(id -g)", "$HOME/.kube/config"},
		true,
	)
}

func init() {
	name := "kube_post_init"
	Plugins[name] = PostInitPlugin{}
}
