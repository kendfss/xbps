package main

import (
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"slices"
	"strings"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
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
	if os.Args[1] == "-d" || os.Args[1] == "--debug" {
		debugInfo()
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
	name := filepath.Base(os.Args[0])
	logf("usage:\n\n\t%s [flags] [subcommand] [subcommand flags] [subcommand args]\n\n", name)
	longest := 0
	for key := range optionTable {
		longest = max(longest, len(key))
	}
	logf("the subcommands are:\n\n")
	for _, key := range slices.Sorted(maps.Keys(optionTable)) {
		val := optionTable[key]
		logf("\t%s%s", key, strings.Repeat(" ", max(longest-len(key)+1, 0)))
		logf("\truns %q%swith desired flags and args\n", val, strings.Repeat(" ", longest+1-len(key)))
	}
	logf("\nthe flags are:\n\n")
	flags := map[string]string{
		"-h, --help":  "print this message",
		"-d, --debug": "print debug info",
	}
	longestFlag := 0
	for key := range flags {
		longestFlag = max(longestFlag, len(key))
	}
	for _, flag := range slices.Sorted(maps.Keys(flags)) {
		text := flags[flag]
		logf("\t%s%s\t%s\n", flag, strings.Repeat(" ", longest-longestFlag), text)
	}
	logf("\n")
}

// debugInfo prints the debug information to stderr
func debugInfo() {
	w, _, _ := termSize()
	block := fmt.Sprintf("version %s\n commit %s\n   date %s\n", version, commit, date)
	info, ok := debug.ReadBuildInfo()
	if ok {
		block += fmt.Sprintf("   home https://%s\n", info.Path)
	}
	widest := 0
	for line := range strings.Lines(block) {
		widest = max(widest, len(line))
	}
	for line := range strings.Lines(block) {
		line = strings.Repeat(" ", max(w/2-widest/2+1, 0)) + line
		logf("%s", line)
	}
}
