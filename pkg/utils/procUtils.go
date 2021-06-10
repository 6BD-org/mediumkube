package utils

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
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

// SameProcInThisContext compare if two processes are
// same
// Under linux, onle first 15 characters are kept
func SameProcInThisContext(p1 string, p2 string) bool {
	prefix := "mediumkube-"
	return strings.Contains(p1, prefix) && strings.Contains(p2, prefix) && p1[:15] == p2[:15]
}
