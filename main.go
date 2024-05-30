package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	command := "commix -r req.txt --ignore-stdin --batch"

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	output := osExec(command, wd)

	fmt.Println(output)
}

func osExec(command, wd string) (output string) {
	var (
		stdout, stderr string
		outbuf, errbuf bytes.Buffer
		exitCode       int
	)

	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Dir = wd
	cmd.Env = os.Environ()

	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Stdin = os.Stdin

	err := cmd.Run()

	stdout = outbuf.String()
	stderr = errbuf.String()

	if err != nil {
		log.Printf("Error executing command: %v", err)
		exitCode = getErrorExitCode(err)
	} else {
		exitCode = cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	}

	if stderr == "" && exitCode != 0 {
		stderr += err.Error()
	}

	output = strings.Join([]string{stdout, stderr, "Exit Status: " + strconv.Itoa(exitCode)}, "\n")
	output = strings.TrimPrefix(output, "\n")
	output = strings.TrimSuffix(output, "\n")

	return output
}

func getErrorExitCode(err error) int {
	if exitError, ok := err.(*exec.ExitError); ok {
		ws := exitError.Sys().(syscall.WaitStatus)
		return ws.ExitStatus()
	}
	return 1
}
