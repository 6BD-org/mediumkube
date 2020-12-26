package services

import (
	"mediumkube/common"
	"mediumkube/configurations"
	"mediumkube/utils"

	"github.com/libvirt/libvirt-go"
)

type LibvirtService struct {
	config *common.OverallConfig
	conn   *libvirt.Connect
}

func init() {
	conn, err := libvirt.NewConnect("qemu:///system")
	utils.CheckErr(err)
	InitLibvritService(
		LibvirtService{
			config: configurations.Config(),
			conn:   conn,
		},
	)
}
