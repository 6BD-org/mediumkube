package mediumssh

import (
	"mediumkube/utils"

	"golang.org/x/crypto/ssh"
)

// SSHLogin get client or die
func SSHLogin(username string, host string, keyPath string) *ssh.Client {
	key := utils.ReadByte(keyPath)
	signer, err := ssh.ParsePrivateKey(key)
	utils.CheckErr(err)

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}
	client, err := ssh.Dial("tcp", host, config)
	return client
}

func init() {

}
