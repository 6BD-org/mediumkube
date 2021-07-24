package mgrpc

import (
	"fmt"
	_mgrpc "mediumkube/daemon/mgrpc"
	"mediumkube/pkg/common"
	"mediumkube/pkg/utils"

	"google.golang.org/grpc"
)

func NewMediumkubeClientOrDie(config *common.OverallConfig, addr string) _mgrpc.DomainSerciceClient {
	// TODO: Security options
	conn, err := grpc.Dial(addr+":"+fmt.Sprint(config.Overlay.GRPCPort), grpc.WithInsecure())
	utils.CheckErr(err)
	return _mgrpc.NewDomainSerciceClient(conn)
}
