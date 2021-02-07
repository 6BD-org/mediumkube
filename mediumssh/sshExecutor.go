package mediumssh

import (
	"bufio"
	"mediumkube/utils"
	"os"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
	"k8s.io/klog/v2"
)

type SSHClient struct {
	client *ssh.Client
}

func _mkdirCmd(dir string) []string {
	return []string{"mkdir", "-p", dir}
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

// Transfer a file from local file system to a ssh server
func (sc SSHClient) Transfer(srcPath string, targetPath string) {
	sess, err := sc.client.NewSession()
	utils.CheckErr(err)

	tgtDir := utils.GetFileDir(targetPath)
	err = sess.Run(strings.Join(_mkdirCmd(tgtDir), " "))
	utils.CheckErr(err)

	wg := sync.WaitGroup{}

	file, err := os.Open(srcPath)
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	pipe, err := sess.StdinPipe()
	utils.CheckErr(err)

	buf := make([]byte, 1024)
	scanner.Buffer(buf, 1024)
	go func() {
		wg.Add(1)
		defer wg.Done()
		sess.Run(strings.Join([]string{"tee", targetPath}, " "))
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		for scanner.Scan() {
			pipe.Write(scanner.Bytes())
		}
		err = pipe.Close()
		klog.Error(err)

	}()

	wg.Wait()

}

func init() {

}
