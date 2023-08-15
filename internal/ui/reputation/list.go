// Package reputation contains UI objects for displaying reputation in various categories
package reputation

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/sot-api/internal/api/structs"
)

// Tree displays reputation categories and their emblem unlock progress
type Tree struct {
	rootNode *node

	tree *widget.Tree

	widget.BaseWidget
}

// NewTree initializes the view's tree and bindings
func NewTree() *Tree {
	t := &Tree{
		rootNode: &node{Label: "Root node placeholder"},
	}
	t.ExtendBaseWidget(t)
	t.createTree()
	return t
}

func (t *Tree) createTree() {
	t.tree = widget.NewTree(
		t.getChildTreeNodes,
		t.isBranch,
		func(branch bool) fyne.CanvasObject {
			if branch {
				return widget.NewLabel("Branch template")
			}
			w := widget.NewLabel("Leaf template")
			w.Wrapping = fyne.TextWrapWord
			return w
		},
		func(id widget.TreeNodeID, _ bool, o fyne.CanvasObject) {
			node := t.rootNode.findNode(id)
			o.(*widget.Label).SetText(node.Label)
		},
	)
}

func (t *Tree) isBranch(id widget.TreeNodeID) bool {
	return !t.rootNode.findNode(id).isLeaf()
}

func (t *Tree) getChildTreeNodes(id widget.TreeNodeID) []widget.TreeNodeID {
	result := t.rootNode.childIDsOf(id)
	return result
}

// CreateRenderer implements the fyne.Widget interface
func (t *Tree) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		t.tree,
	)
}

// SetReputations can be used to update the reputations
// FIXME: displayed order changes on every update
func (t *Tree) SetReputations(data structs.Reputations) {
	t.rootNode = newRootNode(nodeImplReputations(data))
	t.tree.Refresh()
}
