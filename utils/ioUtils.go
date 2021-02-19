package utils

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

// ExecWithStdio execute the command with stdio attached
// so that runtime output can be captured
func ExecWithStdio(cmd *exec.Cmd) (string, error) {

	var stdoutBuf, stderrBuf bytes.Buffer

	/* Error occurred when copying out/err to std */
	var errCpyOut, errCpyErr error

	var err error

	/* IO Pipe of command */
	stdOutIn, _ := cmd.StdoutPipe()
	stdErrIn, _ := cmd.StderrPipe()

	defer stdOutIn.Close()
	defer stdErrIn.Close()

	err = cmd.Start()
	CheckErr(err)

	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	go func() {
		_, errCpyOut = io.Copy(stdout, stdOutIn)

	}()

	go func() {
		_, errCpyErr = io.Copy(stderr, stdErrIn)
	}()

	CheckErr(errCpyOut)
	CheckErr(errCpyErr)
	err = cmd.Wait()
	return string(stdoutBuf.Bytes()), err
}

// AttachAndExec Exec the command with stdio attached
func AttachAndExec(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
