package reputation

import (
	"reflect"
	"testing"
)

type nodeImplBranch struct {
	children []nodeProvider
}

func (n nodeImplBranch) Label() string {
	return "Branch"
}

func (n nodeImplBranch) Children() []nodeProvider {
	return n.children
}

type nodeImplLeaf struct{}

func (n nodeImplLeaf) Label() string {
	return "Leaf"
}

func (n nodeImplLeaf) Children() []nodeProvider {
	return nil
}

func TestNewNode(t *testing.T) {
	nodeImpl := nodeImplBranch{
		children: []nodeProvider{
			nodeImplLeaf{},
		},
	}
	rootNode := newRootNode(nodeImpl)
	if rootNode.Label != "Branch" {
		t.Errorf("rootNode.Label = %s, expected 'Branch'", rootNode.Label)
	}
	leafNode := rootNode.findNode("0")
	if leafNode.Label != "Leaf" {
		t.Errorf("rootNode[0].Label = %s, expected 'Leaf'", leafNode.Label)
	}
}

var rootNode = &node{
	Label: "root",
	children: map[string]*node{
		"0": {
			Label: "Depth 1 - Child 1",
			children: map[string]*node{
				"a": {Label: "Depth 2 - Child 1"},
				"b": {Label: "Depth 2 - Child 2"},
				"c": {Label: "Depth 2 - Child 3"},
			},
		},
		"1": {
			Label: "Depth 1 - Child 2",
			children: map[string]*node{
				"a": {Label: "Depth 2 - Child 4"},
				"b": {Label: "Depth 2 - Child 5"},
				"c": {Label: "Depth 2 - Child 6"},
			},
		},
	},
}

func TestNodeFindNode(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want *node
	}{
		{name: "Depth 0", id: "", want: rootNode},
		{name: "Depth 1", id: "0", want: rootNode.children["0"]},
		{name: "Depth 2", id: "1-a", want: rootNode.children["1"].children["a"]},
		{name: "Nonexistant", id: "1-d", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rootNode.findNode(tt.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("node.findNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeChildIDsOf(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want []string
	}{
		{name: "Depth 0", id: "", want: []string{"0", "1"}},
		{name: "Depth 1", id: "0", want: []string{"0-a", "0-b", "0-c"}},
		{name: "Depth 2", id: "0-a", want: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rootNode.childIDsOf(tt.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("node.childIDsOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
