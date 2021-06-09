package utils

import (
	"regexp"
	"strings"
)

// SplitCmd Split a command into tokens
func SplitCmd(cmd string) []string {
	cmd = strings.TrimSpace(cmd)
	return regexp.MustCompile(`[\s\t]+`).Split(cmd, -1)
}
