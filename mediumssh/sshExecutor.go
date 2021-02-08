package mediumssh

import (
	"bufio"
	"io"
	"mediumkube/utils"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"k8s.io/klog/v2"
)

// SSHClient Mediunkube managed ssh client
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

	if len(targetPath) == 0 {
		klog.Error("Illegal target path")
		return
	}

	if targetPath[len(targetPath)-1] == '/' {
		targetPath = path.Join(targetPath, utils.GetFileName(srcPath))
	}

	tgtDir := utils.GetFileDir(targetPath)
	if len(tgtDir) > 0 {
		err = sess.Run(strings.Join(_mkdirCmd(tgtDir), " "))
		utils.CheckErr(err)
		sess, err = sc.client.NewSession()
		utils.CheckErr(err)
	}

	wg := sync.WaitGroup{}

	file, err := os.Open(srcPath)
	utils.CheckErr(err)
	scanner := bufio.NewScanner(file)

	pipe, err := sess.StdinPipe()
	utils.CheckErr(err)

	wg.Add(2)

	go func() {
		defer wg.Done()
		err = sess.Run(strings.Join([]string{"tee", targetPath}, " "))
		if err != nil {
			klog.Error(err)
			return
		}
	}()

	time.Sleep(1 * time.Second)
	go func() {
		defer wg.Done()
		klog.Info("Sending file")
		for scanner.Scan() {
			_, err := pipe.Write(scanner.Bytes())
			if err != nil && err != io.EOF {
				klog.Error(err)
			}
		}
		scanerr := scanner.Err()
		utils.CheckErr(scanerr)
		if err != nil {
			klog.Error(err)
			return
		}
		sess.Close()

	}()
	wg.Wait()

}

func init() {

}
