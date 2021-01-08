package services

import (
	"mediumkube/common"
	"testing"
)

func TestIpSplit(t *testing.T) {
	a := bridgeSubNet(common.Bridge{
		Inet: "192.168.1.1/24",
	})

	if a != "192.168.1" {
		t.Fail()
	}
}
