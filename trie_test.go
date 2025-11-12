package main

import (
	"strings"
	"testing"
)

var longestCommandLength int

func init() {
	for commandName := range commandTable {
		longestCommandLength = max(longestCommandLength, len(commandName))
	}
}

func TestTrie(t *testing.T) {
	t.Run("methods", func(t *testing.T) {
		t.Run("contains", func(t *testing.T) {
			tests := []struct {
				arg  string
				want bool
			}{
				{"install", true},
				{"fbulk", true},
				{"fetch", true},
				{"reconfigure", true},
				{"remove", true},
				{"rindex", true},
				{"search", false},
				{"rem", false},
				{"reimove", false},
				{"rmove", false},
			}
			for _, test := range tests {
				t.Run(test.arg, func(t *testing.T) {
					have := aliasTrie.contains(test.arg)
					if have != test.want {
						t.Errorf("%q%s have %t, want %t", test.arg, strings.Repeat(" ", longestCommandLength-len(test.arg)), have, test.want)
					}
				})
			}
		})
	})
	t.Run("prefix", func(t *testing.T) {
		tests := []struct {
			arg, want string
		}{
			{"install", "i"},
			{"fbulk", "fb"},
			{"fetch", "fe"},
			{"reconfigure", "rec"},
			{"remove", "rem"},
			{"rindex", "ri"},
		}
		longest := 0
		for _, test := range tests {
			longest = max(longest, len(test.arg))
		}
		for _, test := range tests {
			t.Run(test.arg, func(t *testing.T) {
				have := aliasTrie.alias(test.arg)
				if have != test.want {
					t.Errorf("%s%s have %q, want %q", test.arg, strings.Repeat(" ", longest-len(test.arg)), have, test.want)
				}
			})
		}
	})
}
