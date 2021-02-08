package commands

import (
	"mediumkube/common"
	"path/filepath"
)

// Help print help info. return false if it's not help command
func Help(handler Handler, args []string) bool {
	if len(args) >= 2 && args[1] == "help" {
		handler.Help()
		return true
	}
	return false
}

// logPath where the log in vms are mounted
func logPath(config *common.OverallConfig) string {
	return filepath.Join(config.TmpDir, "vmlogs")
}
