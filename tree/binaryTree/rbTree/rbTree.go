package rbTree

import (
	"algorithms/tree/binaryTree/genericBinaryTree"
)

const (
	black = true
	red   = false
	left  = true
)

type RBT struct {
	genericBinaryTree.GBT
}

func (t *RBT) setColor(node *genericBinaryTree.GBTElement, color bool) {
	node.SideValue = color
}

func (t *RBT) color(node *genericBinaryTree.GBTElement) (black bool) {
	return t.IsNil(node) || node.SideValue.(bool)
}

func (t *RBT) otherSideNode(side bool, node *genericBinaryTree.GBTElement) (*genericBinaryTree.GBTElement) {
	if side == left {
		return node.Right
	} else {
		return node.Left
	}
}
func (t *RBT) invDirRotation(side bool, node *genericBinaryTree.GBTElement) (interface{}) {
	if side == left {
		return t.RightRotate(node)
	} else {
		return t.LeftRotate(node)
	}
}
func (t *RBT) sameSideNode(side bool, node *genericBinaryTree.GBTElement) (*genericBinaryTree.GBTElement) {
	if side == left {
		return node.Left
	} else {
		return node.Right
	}
}
func (t *RBT) sameDirRotation(side bool, node *genericBinaryTree.GBTElement) (interface{}) {
	if side == left {
		return t.LeftRotate(node)
	} else {
		return t.RightRotate(node)
	}
}

func (t *RBT) Insert(node interface{}) (interface{}) {
	n := t.GBT.Insert(node).(*genericBinaryTree.GBTElement)
	t.setColor(n, red)
	t.insertFix(n)
	return n
}

func (t *RBT) insertFix(node interface{}) () {
	n := node.(*genericBinaryTree.GBTElement)
	//only can violate property 3: both left and right children of red node must be black
	for !t.color(n.Parent) && !t.color(n) {
		grandNode := n.Parent.Parent //must be black
		uncleNode := grandNode.Right
		if n.Parent == uncleNode {
			uncleNode = grandNode.Left
		}
		//case1: uncle node is red
		if !t.color(uncleNode) {
			t.setColor(grandNode, red)
			t.setColor(grandNode.Left, black)
			t.setColor(grandNode.Right, black)
			n = grandNode
			//case2&3: uncle node is black
		} else {
			side := n.Parent == grandNode.Left
			t.setColor(grandNode, red)
			//case 2 n is right child of parent
			if n == t.otherSideNode(side, n.Parent) {
				t.sameDirRotation(side, n.Parent)
			}
			//case 3 n is left child of parent
			t.setColor(t.sameSideNode(side, grandNode), black)
			t.invDirRotation(side,grandNode)
		}
	}
	t.setColor(t.Root().(*genericBinaryTree.GBTElement), black)
}

func (t *RBT) Delete(key uint32) (interface{}) {
	deleteNonCompletedNode := func(node *genericBinaryTree.GBTElement) (deletedNode *genericBinaryTree.GBTElement, nextNode *genericBinaryTree.GBTElement) {
		var reConnectedNode *genericBinaryTree.GBTElement
		if t.IsNil(node.Left) {
			reConnectedNode = node.Right
		} else {
			reConnectedNode = node.Left
		}
		//mean's another black color
		reConnectedNode.Parent = node.Parent
		if t.IsNil(node.Parent) {
			t.NilNode.Left = reConnectedNode
			t.NilNode.Right = reConnectedNode
		} else if node.Parent.Right == node {
			node.Parent.Right = reConnectedNode
		} else {
			node.Parent.Left = reConnectedNode
		}
		return node, reConnectedNode
	}
	node := t.Search(key).(*genericBinaryTree.GBTElement)
	if t.IsNil(node) {
		return node
	}
	var deletedNode, reConnectedNode *genericBinaryTree.GBTElement
	if t.IsNil(node.Left) || t.IsNil(node.Right) {
		deletedNode, reConnectedNode = deleteNonCompletedNode(node)
	} else {
		successor := t.Successor(node, t.Root()).(*genericBinaryTree.GBTElement)
		_key, _value := successor.Key, successor.Value
		node.Key, node.Value = _key, _value
		deletedNode, reConnectedNode = deleteNonCompletedNode(successor)
	}
	if t.color(deletedNode) {
		//Now, reConnectedNode is black-black or black-red
		t.deleteFix(reConnectedNode)
	}
	//recover NilNode
	t.NilNode.Parent = t.NilNode
	return node
}

func (t *RBT) deleteFix(node interface{}) () {
	n := node.(*genericBinaryTree.GBTElement)
	brotherNode := func(side bool, node *genericBinaryTree.GBTElement) (*genericBinaryTree.GBTElement) {
		return t.otherSideNode(side, node.Parent)
	}
	case1 := func(side bool, node *genericBinaryTree.GBTElement) {
		t.setColor(node.Parent, red)
		t.setColor(brotherNode(side, node), black)
		t.sameDirRotation(side, node.Parent)
	}
	case2 := func(side bool, node *genericBinaryTree.GBTElement) (*genericBinaryTree.GBTElement) {
		t.setColor(brotherNode(side, node), red)
		return node.Parent
	}
	case3 := func(side bool, node *genericBinaryTree.GBTElement) {
		brother := brotherNode(side, node)
		t.setColor(brother, red)
		t.setColor(t.sameSideNode(side, brother), black)
		t.invDirRotation(side, brother)
	}
	case4 := func(side bool, node *genericBinaryTree.GBTElement) (*genericBinaryTree.GBTElement) {
		brother := brotherNode(side, node)
		t.setColor(brother, t.color(node.Parent))
		t.setColor(node.Parent, black)
		t.setColor(t.otherSideNode(side, brother), black)
		t.sameDirRotation(side, node.Parent)
		return t.Root().(*genericBinaryTree.GBTElement)
	}
	//n always points to the black-black or black-red node.The purpose is to remove the additional black color,
	//which means add a black color in the same side or reduce a black color in the other side
	for n != t.Root() && t.color(n) {
		side := n == n.Parent.Left
		brother := brotherNode(side, n)
		//case 1 brother node is red, so parent must be black.Turn brother node to a black one, convert to case 2,3,4
		if !t.color(brother) {
			case1(side, n)
			//case 2, 3, 4 brother node is black
		} else {
			//case 2 move black-blcak or blcak-red node up
			if t.color(brother.Left) && t.color(brother.Right) {
				n = case2(side, n)
				//case 3 convert to case 4
			} else if t.color(t.otherSideNode(side, brother)) {
				case3(side, n)
				//case 4 add a black to left, turn black-black or black-red to black or red
			} else {
				n = case4(side, n)
			}
		}

	}
	t.setColor(n, black)

}

func New() *RBT {
	t := new(RBT)
	t.Init()
	t.GBT.Object = t
	return t
}