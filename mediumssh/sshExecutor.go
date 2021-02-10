package mediumssh

import (
	"bufio"
	"fmt"
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

// SSHMeta stores extra metadata used my mediumssh
type sshMeta struct {
	username string
}

// SSHClient Mediunkube managed ssh client
// I suggest that the only way to obtain SSHClient is using SSHLogin function
type SSHClient struct {
	client *ssh.Client
	meta   *sshMeta
}

func _mkdirCmd(dir string) []string {
	return []string{"mkdir", "-p", dir}
}

func _chownCmd(dir string, user string, group string) []string {
	return []string{"chown", fmt.Sprintf("%v:%v", user, group), "-R", dir}

}
func _sudo(cmd []string, sudo bool) []string {
	if sudo {
		return append([]string{"sudo"}, cmd...)
	}
	return cmd

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
	meta := sshMeta{username: username}
	return SSHClient{client: client, meta: &meta}
}

// AttachAndExecute Attach to the session and Execute a command
// Under this mode, input of stdio can be captured by remote host
func (sc SSHClient) AttachAndExecute(cmd []string, sudo bool) {
	cmd = _sudo(cmd, sudo)
	sess, err := sc.client.NewSession()
	utils.CheckErr(err)
	sess.Stderr = os.Stderr
	sess.Stdin = os.Stdin
	sess.Stdout = os.Stdout

	err = sess.Run(strings.Join(cmd, " "))
	utils.CheckErr(err)
}

// Execute a command remotely. The output and err will be printed on current Stdio
func (sc SSHClient) Execute(cmd []string, sudo bool) {
	cmd = _sudo(cmd, sudo)
	sess, err := sc.client.NewSession()
	utils.CheckErr(err)
	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr
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
		// Make dir as root
		sc.Execute(_mkdirCmd(tgtDir), true)
		// Recover the ownership of dir
		user := sc.meta.username
		group := user
		sc.Execute(_chownCmd(tgtDir, user, group), true)

	}

	wg := sync.WaitGroup{}

	file, err := os.Open(srcPath)
	utils.CheckErr(err)
	scanner := bufio.NewScanner(file)

	pipe, err := sess.StdinPipe()
	utils.CheckErr(err)

	wg.Add(2)

	go func() {
		// Start a process in remote host that redirect stdin to a file
		defer wg.Done()
		err = sess.Run(strings.Join([]string{"tee", targetPath}, " "))
		if err != nil {
			if !strings.Contains(err.Error(), "signal PIPE") {
				klog.Error(err)
				return
			}
		}
	}()

	// Wait for remote process to start
	time.Sleep(1 * time.Second)

	go func() {
		// Start a local process that read data from source file
		// and write to the pipeline
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
