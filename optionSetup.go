package main

import (
	"os"
	"path/filepath"
	"strings"
)

var optionTable map[string]string

func init() {
	var err error
	optionTable, err = children("xbps")
	if err != nil {
		panic(err)
	}
}

func children(name string) (map[string]string, error) {
	roots := filepath.SplitList(os.Getenv("PATH"))
	table := map[string]string{}
	for _, root := range roots {
		matches, err := filepath.Glob(filepath.Join(root, name+"-*"))
		if err != nil {
			return nil, err
		}
		for _, match := range matches {
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
