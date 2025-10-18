package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func handleStdin(stdin io.WriteCloser) {
	defer stdin.Close()

	// Write to command's stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "exit" {
			break
		}
		fmt.Fprintln(stdin, line)
	}
}

func handleStdout(stdout io.Reader) {
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func handleStderr(stderr io.Reader) {
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		logf("%s\n", scanner.Text())
	}
}

func run(cmd *exec.Cmd) {
	// Get pipes for stdin, stdout, stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		logf("Error creating stdin pipe: %v\n", err)
		os.Exit(1)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logf("Error creating stdout pipe: %v\n", err)
		os.Exit(1)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logf("Error creating stderr pipe: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	// Handle I/O in goroutines
	go handleStdin(stdin)
	go handleStdout(stdout)
	go handleStderr(stderr)

	if err := cmd.Wait(); err != nil {
		logf("%s\n", err)
		os.Exit(err.(*exec.ExitError).ExitCode())
	}
}
