package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	commandTable map[string]string // table of commands and their full names
	aliasTable   map[string]string // table of aliases and their respective commands
	aliasTrie    trie              // computes the shortest unique prefix of each command to populate the alias table
)

func init() {
	var err error
	commandTable, err = children("xbps", os.Getenv("PATH"))
	if err != nil {
		logf("couldn't find options: %s", err)
		os.Exit(1)
	}
	for command := range commandTable {
		aliasTrie.learn(command)
	}
	aliasTable = make(map[string]string, len(commandTable))
	for key := range commandTable {
		aliasTable[aliasTrie.alias(key)] = key
	}
}

// children discovers the names of all child commands of the given base command available on the PATH
func children(name, PATH string) (map[string]string, error) {
	roots := filepath.SplitList(PATH)
	table := map[string]string{}
	for _, root := range roots {
		matches, err := filepath.Glob(filepath.Join(root, name+"-*"))
		if err != nil {
			return nil, err
		}
		for _, match := range matches {
			if !isExecutable(match) {
				continue
			}
			base := filepath.Base(match)
			_, found := table[base]
			if found {
				continue
			}
			key := strings.Join(strings.Split(base, "-")[1:], "-")
			table[key] = base
		}
	}
	return table, nil
}

// isExecutable checks if the file at the given path is isExecutable
func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if runtime.GOOS == "windows" {
		// this part is technically irrelevant, but in the name of future reference (ie who of the void really knows what WSL is doing these days):
		_, found := map[string]bool{
			".exe": true, ".com": true, ".bat": true,
			".cmd": true, ".ps1": true,
		}[strings.ToLower(filepath.Ext(path))]
		return found
	}
	// On Unix-like systems, check executable bits
	return info.Mode().Perm()&0o111 != 0
}
