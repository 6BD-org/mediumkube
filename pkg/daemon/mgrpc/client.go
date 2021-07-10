package mgrpc

import (
	"fmt"
	"mediumkube/pkg/common"
	"mediumkube/pkg/utils"

	"google.golang.org/grpc"
)

func NewMediumkubeClientOrDie(config *common.OverallConfig, addr string) DomainSerciceClient {
	// TODO: Security options
	conn, err := grpc.Dial(addr+":"+fmt.Sprint(config.Overlay.GRPCPort), grpc.WithInsecure())
	utils.CheckErr(err)
	return NewDomainSerciceClient(conn)
}
