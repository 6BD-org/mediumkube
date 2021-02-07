package mediumssh

import (
	"mediumkube/utils"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	client *ssh.Client
}

// SSHLogin get client or die
func SSHLogin(username string, host string, keyPath string) SSHClient {
	key := utils.ReadByte(keyPath)
	signer, err := ssh.ParsePrivateKey(key)
	utils.CheckErr(err)

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", host, config)
	utils.CheckErr(err)
	return SSHClient{client: client}
}

// Execute Execute a command
func (sc SSHClient) Execute(cmd []string, sudo bool) {
	sess, err := sc.client.NewSession()
	utils.CheckErr(err)
	sess.Stderr = os.Stderr
	sess.Stdin = os.Stdin
	sess.Stdout = os.Stdout

	err = sess.Run(strings.Join(cmd, " "))
	utils.CheckErr(err)
}

func init() {

}
