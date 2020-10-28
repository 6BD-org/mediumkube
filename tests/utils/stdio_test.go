package utils

import (
	"mediumkube/utils"
	"os/exec"
	"testing"
)

func TestStdio(t *testing.T) {
	t.Log("Doing stdio test. Please check command line output. If files in this dir are printed, then the test is passed")
	t.Log("This is not a proper way to do ut though :(")
	utils.ExecWithStdio(exec.Command("ls"))
}
