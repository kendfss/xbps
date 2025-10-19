package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var optionTable map[string]string // table of options and their full names

func init() {
	var err error
	optionTable, err = children("xbps", os.Getenv("PATH"))
	if err != nil {
		logf("couldn't find options: %s", err)
		os.Exit(1)
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
			if !executable(match) {
				continue
			}
			base := filepath.Base(match)
			_, found := table[base]
			if found {
				continue
			}
			key := strings.Split(base, "-")[1]
			table[key] = base
		}
	}
	return table, nil
}

// executable checks if the file at the given path is executable
func executable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if runtime.GOOS == "windows" {
		_, found := map[string]bool{
			".exe": true, ".com": true, ".bat": true,
			".cmd": true, ".ps1": true,
		}[filepath.Ext(path)]
		return found
	}
	// On Unix-like systems, check executable bits
	return info.Mode().Perm()&0111 != 0
}
