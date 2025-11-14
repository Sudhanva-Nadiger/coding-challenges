package huffmantree

import "container/heap"

type Node struct {
	weight int64
	char   rune
	left   *Node
	right  *Node
}

type HuffManTree struct {
	root         *Node
	FrequencyMap *map[rune]int64
}

func NewHuffManTree(frequencyMap map[rune]int64) *HuffManTree {
	return &HuffManTree{
		root:         nil,
		FrequencyMap: &frequencyMap,
	}
}

func (ht *HuffManTree) BuildHuffManTree() {
	pq := BuildPriorityQueue(ht.FrequencyMap)

	for pq.Len() > 1 {
		node1 := heap.Pop(pq).(*Node)
		node2 := heap.Pop(pq).(*Node)

		sum := node1.weight + node2.weight

		mergedNode := MergeTree(sum, node1, node2)

		heap.Push(pq, mergedNode)
	}

	ht.root = heap.Pop(pq).(*Node)
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
