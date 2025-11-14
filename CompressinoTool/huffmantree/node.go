package huffmantree

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
