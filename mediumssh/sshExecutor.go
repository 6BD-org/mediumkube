package mediumssh

import (
	"bufio"
	"fmt"
	"io"
	"mediumkube/utils"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
	"k8s.io/klog/v2"
)

var (
	sigtable map[os.Signal]byte = make(map[os.Signal]byte)
)

/*

Opcode 	Argument 	Description 	Reference 	Note
0	TTY_OP_END	Indicates end of options.	[RFC4250]
1	VINTR	Interrupt character; 255 if none. Similarly for the other characters. Not all of these characters are supported on all systems.	[RFC4254]	Section 8
2	VQUIT	The quit character (sends SIGQUIT signal on POSIX systems).	[RFC4254]	Section 8
3	VERASE	Erase the character to left of the cursor.	[RFC4254]	Section 8
4	VKILL	Kill the current input line.	[RFC4254]	Section 8
5	VEOF	End-of-file character (sends EOF from the terminal).	[RFC4254]	Section 8
6	VEOL	End-of-line character in addition to carriage return and/or linefeed.	[RFC4254]	Section 8
7	VEOL2	Additional end-of-line character.	[RFC4254]	Section 8
8	VSTART	Continues paused output (normally control-Q).	[RFC4254]	Section 8
9	VSTOP	Pauses output (normally control-S).	[RFC4254]	Section 8
10	VSUSP	Suspends the current program.	[RFC4254]	Section 8
11	VDSUSP	Another suspend character.	[RFC4254]	Section 8
12	VREPRINT	Reprints the current input line.	[RFC4254]	Section 8
13	VWERASE	Erases a word left of cursor.	[RFC4254]	Section 8
14	VLNEXT	Enter the next character typed literally, even if it is a special character	[RFC4254]	Section 8
15	VFLUSH	Character to flush output.	[RFC4254]	Section 8
16	VSWTCH	Switch to a different shell layer.	[RFC4254]	Section 8
17	VSTATUS	Prints system status line (load, command, pid, etc).	[RFC4254]	Section 8
18	VDISCARD	Toggles the flushing of terminal output.	[RFC4254]	Section 8
19-29	Unassigned
30	IGNPAR	The ignore parity flag. The parameter SHOULD be 0 if this flag is FALSE, and 1 if it is TRUE.	[RFC4254]	Section 8
31	PARMRK	Mark parity and framing errors.	[RFC4254]	Section 8
32	INPCK	Enable checking of parity errors.	[RFC4254]	Section 8
33	ISTRIP	Strip 8th bit off characters.	[RFC4254]	Section 8
34	INLCR	Map NL into CR on input.	[RFC4254]	Section 8
35	IGNCR	Ignore CR on input.	[RFC4254]	Section 8
36	ICRNL	Map CR to NL on input.	[RFC4254]	Section 8
37	IUCLC	Translate uppercase characters to lowercase.	[RFC4254]	Section 8
38	IXON	Enable output flow control.	[RFC4254]	Section 8
39	IXANY	Any char will restart after stop.	[RFC4254]	Section 8
40	IXOFF	Enable input flow control.	[RFC4254]	Section 8
41	IMAXBEL	Ring bell on input queue full.	[RFC4254]	Section 8
42	IUTF8	Terminal input and output is assumed to be encoded in UTF-8.	[RFC8160]
43-49	Unassigned
50	ISIG	Enable signals INTR, QUIT, [D]SUSP.	[RFC4254]	Section 8
51	ICANON	Canonicalize input lines.	[RFC4254]	Section 8
52	XCASE	Enable input and output of uppercase characters by preceding their lowercase equivalents with "\".	[RFC4254]	Section 8
53	ECHO	Enable echoing.	[RFC4254]	Section 8
54	ECHOE	Visually erase chars.	[RFC4254]	Section 8
55	ECHOK	Kill character discards current line.	[RFC4254]	Section 8
56	ECHONL	Echo NL even if ECHO is off.	[RFC4254]	Section 8
57	NOFLSH	Don't flush after interrupt.	[RFC4254]	Section 8
58	TOSTOP	Stop background jobs from output.	[RFC4254]	Section 8
59	IEXTEN	Enable extensions.	[RFC4254]	Section 8
60	ECHOCTL	Echo control characters as ^(Char).	[RFC4254]	Section 8
61	ECHOKE	Visual erase for line kill.	[RFC4254]	Section 8
62	PENDIN	Retype pending input.	[RFC4254]	Section 8
63-69	Unassigned
70	OPOST	Enable output processing.	[RFC4254]	Section 8
71	OLCUC	Convert lowercase to uppercase.	[RFC4254]	Section 8
72	ONLCR	Map NL to CR-NL.	[RFC4254]	Section 8
73	OCRNL	Translate carriage return to newline (output).	[RFC4254]	Section 8
74	ONOCR	Translate newline to carriage return-newline (output).	[RFC4254]	Section 8
75	ONLRET	Newline performs a carriage return (output).	[RFC4254]	Section 8
76-89	Unassigned
90	CS7	7 bit mode.	[RFC4254]	Section 8
91	CS8	8 bit mode.	[RFC4254]	Section 8
92	PARENB	Parity enable.	[RFC4254]	Section 8
93	PARODD	Odd parity, else even.	[RFC4254]	Section 8
94-127	Unassigned
128	TTY_OP_ISPEED	Specifies the input baud rate in bits per second.	[RFC4254]	Section 8
129	TTY_OP_OSPEED	Specifies the output baud rate in bits per second.	[RFC4254]	Section 8
130-255	Unassigned

*/

// SSHMeta stores extra metadata used my mediumssh
type sshMeta struct {
	username string
	host     string
	port     int
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
func SSHLogin(username string, host string, keyPath string) *SSHClient {
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
	hostAndPort := strings.Split(host, ":")
	port, _ := strconv.Atoi(hostAndPort[1])
	meta := sshMeta{username: username, host: hostAndPort[0], port: port}
	return &SSHClient{client: client, meta: &meta}
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
// This is done by reading file and send to vm's input pipe,
// and redirect stdin to file using tee
func (sc SSHClient) Transfer(srcPath string, targetPath string) {

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

	sess, err := sc.client.NewSession()
	utils.CheckErr(err)
	pipe, err := sess.StdinPipe()
	utils.CheckErr(err)
	writer := bufio.NewWriter(pipe)

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

	go func() {
		// Start a local process that read data from source file
		// and write to the pipeline
		defer wg.Done()
		buffer := make([]byte, 1024*4)
		for {
			n, err := file.Read(buffer)
			if n > 0 {
				_, errW := writer.Write(buffer[:n])
				utils.CheckErr(errW)
			}
			if err != nil {
				if err == io.EOF {
					break
				} else {
					utils.CheckErr(err)
				}
			}
			utils.CheckErr(err)
		}
		if err != nil {
			klog.Error(err)
			return
		}
		writer.Flush()
		pipe.Close()
	}()
	wg.Wait()

}

// Receive a file from src in vm and save as tgt on host machine
func (sc SSHClient) Receive(src string, tgt string) {
	tgtDir := utils.GetFileDir(tgt)
	if len(tgtDir) > 0 {
		mkdirCmd := _mkdirCmd(tgtDir)
		utils.ExecWithStdio(exec.Command(mkdirCmd[0], mkdirCmd[1:]...))
	}

	wg := sync.WaitGroup{}
	wg.Add(2) // receiver

	sess, err := sc.client.NewSession()
	utils.CheckErr(err)

	pipe, err := sess.StdoutPipe()

	go func() {
		defer wg.Done()
		file, err := os.OpenFile(tgt, os.O_CREATE|os.O_RDWR, os.FileMode(0755))
		defer file.Close()
		writer := bufio.NewWriter(file)

		if err != nil {
			klog.Error(err)
			return
		}

		buffer := make([]byte, 1024*4)
		for {

			n, err := pipe.Read(buffer)
			if n > 0 {
				_, err := writer.Write(buffer[:n])
				utils.CheckErr(err)
				err = writer.Flush()
				utils.CheckErr(err)
			}
			if err != nil {
				if err == io.EOF {
					break
				} else {
					utils.CheckErr(err)
				}
			}

		}
	}()

	go func() {
		defer wg.Done()
		err := sess.Run(strings.Join([]string{"cat", src}, " "))
		if err != nil {
			klog.Error(err)
			return
		}
	}()

	wg.Wait()
}

// Shell launch a shell
func (sc SSHClient) Shell() {
	utils.AttachAndExec(exec.Command("ssh", fmt.Sprintf("%v@%v", sc.meta.username, sc.meta.host)))
}

func init() {
	sigtable[os.Kill] = 0x03
	sigtable[os.Interrupt] = 0x03
}
