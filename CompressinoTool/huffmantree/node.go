package huffmantree

import "strconv"

type Node struct {
	weight int64
	char   rune
	left   *Node
	right  *Node
}

func NewLeafNode(w int64, c rune) *Node {
	return &Node{
		weight: w,
		char:   c,
		left:   nil,
		right:  nil,
	}
}

func NewInternalNode(w int64) *Node {
	return &Node{
		weight: w,
		left:   nil,
		right:  nil,
	}
}

func MergeTree(w int64, l, r *Node) *Node {
	node := NewInternalNode(w)

	node.left = l
	node.right = r

	return node
}

func (n *Node) IsLeadNode() bool {
	return n.char != 0
}

func (root *Node) PrintTree() {
	if root == nil {
		println("Tree is empty")
		return
	}
	printNode(root, "", true)
}

func printNode(node *Node, prefix string, isTail bool) {
	if node == nil {
		return
	}

	// Print current node
	connector := "└── "
	if !isTail {
		connector = "├── "
	}

	// Format node information
	var nodeInfo string
	if node.IsLeadNode() {
		// Leaf node - show character and weight
		charDisplay := string(node.char)
		switch node.char {
		case '\n':
			charDisplay = "\\n"
		case '\t':
			charDisplay = "\\t"
		case ' ':
			charDisplay = "SPACE"
		}
		nodeInfo = charDisplay + " (" + formatInt64(node.weight) + ")"
	} else {
		// Internal node - show only weight
		nodeInfo = "* (" + formatInt64(node.weight) + ")"
	}

	println(prefix + connector + nodeInfo)

	// Prepare prefix for children
	extension := "    "
	if !isTail {
		extension = "│   "
	}
	newPrefix := prefix + extension
	if node.right != nil || node.left != nil {
		printNode(node.left, newPrefix, node.right == nil)
		printNode(node.right, newPrefix, true)
	}
}

func formatInt64(n int64) string {
	return strconv.FormatInt(n, 10)
}
