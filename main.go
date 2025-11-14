package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
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
		help(nil)
		os.Exit(1)
	}
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		help(flag.ErrHelp)
		os.Exit(0)
	}
	if os.Args[1] == "-d" || os.Args[1] == "--debug" {
		debugInfo()
		os.Exit(0)
	}

	commandName, found := commandTable[os.Args[1]]
	if !found {
		childName, found := aliasTable[os.Args[1]]
		if !found {
			logf("unsupported command: %q\n", os.Args[1])
			os.Exit(1)
		}
		commandName = commandTable[childName]
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
func help(err error) {
	block := new(bytes.Buffer)
	name := filepath.Base(os.Args[0])
	fmt.Fprintf(block, "usage:\n\n\t%s [flags] [command/alias] [command flags] [command args]\n\n", name)
	longestAlias, longestCommand := 0, 0
	for alias, command := range aliasTable {
		longestCommand = max(longestCommand, len(command))
		longestAlias = max(longestAlias, len(alias))
	}
	fmt.Fprint(block, "the subcommands are:\n\n")
	for _, alias := range slices.Sorted(maps.Keys(aliasTable)) {
		command := aliasTable[alias]
		full := commandTable[command]
		fmt.Fprintf(block, "\t%s%s", command, strings.Repeat(" ", max(longestCommand-len(command)+1, 0)))
		fmt.Fprintf(block, "to run %q%swith desired flags and args\t(aliased to %q)\n", full, strings.Repeat(" ", longestCommand+1-len(command)), alias)
	}
	fmt.Fprint(block, "\nthe flags are:\n\n")
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
		fmt.Fprintf(block, "\t%s%s\t %s\n", flag, strings.Repeat(" ", longestCommand-longestFlag), text)
	}
	dst := os.Stderr
	if err != nil {
		dst = os.Stdout
	}
	io.Copy(dst, block)
}

// debugInfo prints the debug information to stderr
func debugInfo() {
	w, _, _ := termSize()
	block := new(strings.Builder)
	fmt.Fprintf(block, "version %s\n commit %s\n   date %s\n", version, commit, date)
	info, ok := debug.ReadBuildInfo()
	if ok {
		fmt.Fprintf(block, "   home https://%s", info.Path)
	}
	widest := 0
	for line := range strings.Lines(block.String()) {
		widest = max(widest, len(line))
	}
	content := block.String()
	block.Reset()
	for line := range strings.Lines(content) {
		line = strings.Repeat(" ", max(w/2-widest/2+1, 0)) + line
		fmt.Fprintf(block, "%s", line)
	}
	fmt.Println(block.String())
}
