package huffmantree

import "container/heap"

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
