package main

import (
	"fmt"
	"maps"
	"os"
	"os/exec"
	"slices"
)

func main() {
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		logf("how to use %q command:\n", os.Args[0])
		for _, key := range slices.Sorted(maps.Keys(optionTable)) {
			val := optionTable[key]
			logf("\t%s %s [flags] [args]: run %q with select flags and args\n", os.Args[0], key, val)
		}
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
	cmd := exec.Command(commandName, argsToPass...)
	run(cmd)
}

func logf(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg, args...)
}
