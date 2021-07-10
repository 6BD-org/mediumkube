package commands

import (
	"context"
	"flag"
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/daemon/mgrpc"
	"mediumkube/pkg/utils"
	"os"

	"k8s.io/klog/v2"
)

type DeployHandler struct {
	flagset *flag.FlagSet
}

func (handler DeployHandler) Desc() string {
	return "Deploy a new K8s cluster"
}

func (handler DeployHandler) Help() {
	handler.flagset.Usage()

}

func (handler DeployHandler) Handle(args []string) {

	config := configurations.Config()

	var name, cpu, mem, disk string
	handler.flagset.StringVar(&name, "name", "", "Name of the domain. Must be cluster-wise unique")
	handler.flagset.StringVar(&cpu, "cpu", "2", "Number of cpu")
	handler.flagset.StringVar(&mem, "memory", "2G", "Size of memory")
	handler.flagset.StringVar(&disk, "disk", "20G", "size of disk")
	handler.flagset.Parse(args[1:])
	if Help(handler, args) {
		handler.Help()
		return
	}
	if name == "" {
		panic("Invalid name")
	}
	// Mediumkube only supports single node currently.
	// Let scheduler to select node in the future
	client := mgrpc.NewMediumkubeClientOrDie(config, config.Overlay.Master)
	configs := make([]*mgrpc.DomainConfig, 0)
	configs = append(configs, &mgrpc.DomainConfig{Cpu: cpu, Memory: mem, Disk: disk, Name: name})
	// TODO: Handler creation over to scheduler
	stream, err := client.DeployDomain(context.TODO(),
		&mgrpc.DomainCreationParam{Config: configs},
	)
	utils.CheckErr(err)
	for {
		resp, err := stream.Recv()
		if err != nil {
			klog.Error(err)
			return
		}
		os.Stdout.Write([]byte(resp.Content))
	}
}

func init() {
	var name = "deploy"
	CMD[name] = DeployHandler{
		flagset: flag.NewFlagSet(name, flag.ExitOnError),
	}
}
