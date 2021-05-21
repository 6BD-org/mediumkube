package utils

import (
	"mediumkube/pkg/utils"
	"reflect"
	"testing"
)

func TestCommandSplit(t *testing.T) {
	cmd := utils.SplitCmd("ls	  -al  	")
	if !reflect.DeepEqual([]string{"ls", "-al"}, cmd) {
		t.Fail()
	}
}
