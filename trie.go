package main

import (
	"fmt"
	"slices"
	"strings"
)

type (
	trie []*node
	node struct {
		parent   *node
		children []*node
		char     rune
		count    int
	}
)

// append creates a new node for the given char
// or increments the relevant counter if a node already exists
func (t *trie) append(char rune) *node {
	for _, branch := range *t {
		if branch.char == char {
			branch.count++
			return branch
		}
	}
	*t = append(*t, &node{nil, nil, char, 1})
	return (*t)[len(*t)-1]
}

// contains checks if the sequence of nodes traversed by a given string
// is contained within a trie and is also represents a terminal trek
// through its nodes. it rejects empty strings
func (t trie) contains(str string) bool {
	var n *node
outer:
	for _, char := range str {
		for _, node := range t {
			if node.char == char {
				n = node
				t = node.children
				continue outer
			}
		}
		return false
	}
	return n != nil && len(n.children) == 0
}

// alias computes the shortest unique alias of a given string
// returns "" if str is empty
// returns "" if there is no matching substring
// returns "" if the longest matching string is a substring of the argument
func (t *trie) alias(str string) string {
	if len(str) == 0 {
		return ""
	}
	i := slices.IndexFunc(*t, holds(rune(str[0])))
	if i < 0 {
		return ""
	}
	tail := (*t)[i]
	if tail.count == 1 {
		return string(tail.char)
	}
outer:
	for _, char := range str[1:] {
		for _, child := range tail.children {
			if child.char == char {
				tail = child
				if child.count == 1 {
					break outer
				}
				break
			}
		}
	}
	out := ""
	for tail != nil {
		out = string(tail.char) + out
		tail = tail.parent
	}
	return out
}

// holds returns a closure which checks nodes for chars equivalent to the initial character
func holds(char rune) func(*node) bool {
	return func(n *node) bool {
		return n.char == char
	}
}

// String represents the trie for command line debugging
func (t trie) String() string {
	nodes := make([]*node, len(t))
	copy(nodes, t)
	name := strings.Split(fmt.Sprintf("%T", t), ".")[1]
	return fmt.Sprintf("%s%s", name, nodes)
}

// String represents a node for command line debugging
func (n node) String() string {
	childChars := make([]rune, len(n.children))
	for i, child := range n.children {
		childChars[i] = child.char
	}
	if n.parent != nil {
		return fmt.Sprintf("%T{%c, %c, %c, %d}", n, n.parent.char, childChars, n.char, n.count)
	}
	name := strings.Split(fmt.Sprintf("%T", n), ".")[1]
	return fmt.Sprintf("%s{nil, %c, %c, %d}", name, childChars, n.char, n.count)
}

// append creates a new child for the given char
// or increments the relevant counter if there is already a child for the char
func (n *node) append(char rune) *node {
	for _, child := range n.children {
		if child.char == char {
			child.count++
			return child
		}
	}
	n.children = append(n.children, &node{n, nil, char, 1})
	return n.children[len(n.children)-1]
}

// learn feeds a new sequence of characters into the trie
func (t *trie) learn(str string) {
	if len(str) == 0 {
		return
	}
	ptr := t.append(rune(str[0]))
	for _, char := range str[1:] {
		ptr = ptr.append(char)
	}
}
