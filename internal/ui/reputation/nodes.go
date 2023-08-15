package reputation

import (
	"slices"
	"strconv"
	"strings"
)

type node struct {
	Label    string
	children map[string]*node
}

type nodeProvider interface {
	// Label returns the label for this node
	Label() string
	// Children returns the children of this node.
	// Return nil or an empty slice to indicate that this node is a leaf.
	Children() []nodeProvider
}

func newRootNode(src nodeProvider) (n *node) {
	n = &node{Label: src.Label()}
	srcChildren := src.Children()
	if srcChildren == nil {
		return
	}
	children := make(map[string]*node, len(srcChildren))
	for i, v := range srcChildren {
		children[strconv.Itoa(i)] = newRootNode(v)
	}
	n.children = children
	return
}

const nodeIDSeparator = "-"

func (n *node) childIDsOf(id string) []string {
	target := n.findNode(id)
	if n.isLeaf() {
		return []string{}
	}
	childIDs := make([]string, len(target.children))
	i := 0
	for k := range target.children {
		if id == "" {
			childIDs[i] = k
		} else {
			childIDs[i] = id + nodeIDSeparator + k
		}
		i++
	}
	slices.Sort(childIDs)
	return childIDs
}

func (n *node) findNode(id string) *node {
	if id == "" {
		return n
	}
	subIDs := strings.SplitN(id, nodeIDSeparator, 2)
	next, ok := n.children[subIDs[0]]
	if !ok {
		return nil
	}
	if len(subIDs) == 1 {
		return next
	}
	return next.findNode(subIDs[1])
}

func (n *node) isLeaf() bool {
	return n.children == nil
}
