package mediumssh

import (
	"mediumkube/pkg/utils"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SFTPWalk Walks the path on given ssh host
func SFTPWalk(sshClient *ssh.Client, path string) []string {
	sftpClient, err := sftp.NewClient(sshClient)
	utils.CheckErr(err)
	walker := sftpClient.Walk(path)
	res := make([]string, 0)
	for walker.Step() {
		if !walker.Stat().IsDir() {
			res = append(res, walker.Path())
		}
	}
	return res
}
