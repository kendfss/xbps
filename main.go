package main

import (
	"fmt"
	"maps"
	"os"
	"os/exec"
	"slices"
)

func main() {
	if len(os.Args) == 1 {
		help()
		os.Exit(1)
	}
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		help()
		os.Exit(0)
	}

	commandName, found := optionTable[os.Args[1]]
	if !found {
		logf("unsupported command: %q\n", os.Args[1])
		os.Exit(1)
	}
	argsToPass := []string{}
	if len(os.Args) > 2 {
		argsToPass = os.Args[2:]
	}

	run(exec.Command(commandName, argsToPass...))
}

// logf prints a message to stderr
func logf(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg, args...)
}

// run configures io and executes the desired command
func run(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logf("%s\n", err)
		if err, ok := err.(*exec.ExitError); ok {
			os.Exit(err.ExitCode())
		}
		os.Exit(1)
	}
}

// help prints the usage information to stderr
func help() {
	logf("how to use %q command:\n", os.Args[0])
	for _, key := range slices.Sorted(maps.Keys(optionTable)) {
		val := optionTable[key]
		logf("\t%s %s [flags] [args]: run %q with desired flags and args\n", os.Args[0], key, val)
	}
}
