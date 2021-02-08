package utils

import (
	"fmt"
	"io/ioutil"
	"path"
)

// GetLinuxProcCmdOrEmpty get command of process, otherwise, return empty string
func GetLinuxProcCmdOrEmpty(pid int) string {
	cmdPath := path.Join("/proc", fmt.Sprintf("%v", pid), "cmdline")
	bt, err := ioutil.ReadFile(cmdPath)
	if err == nil {
		return string(bt)
	}
	return ""
}
