package huffmantree

import (
	"container/heap"
	"sort"
)

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	// First compare by weight
	if pq[i].weight != pq[j].weight {
		return pq[i].weight < pq[j].weight
	}

	// If weights are equal, use tie-breaker for deterministic ordering
	// For leaf nodes, compare by character value
	if pq[i].isLeadNode() && pq[j].isLeadNode() {
		return pq[i].char < pq[j].char
	}

	if pq[i].isLeadNode() {
		return true
	}
	if pq[j].isLeadNode() {
		return false
	}

	return getMinChar(pq[i]) < getMinChar(pq[j])
}

func getMinChar(node *Node) rune {
	if node.isLeadNode() {
		return node.char
	}

	minChar := rune(0x10FFFF) // Max unicode value
	if node.left != nil {
		leftMin := getMinChar(node.left)
		if leftMin < minChar {
			minChar = leftMin
		}
	}
	if node.right != nil {
		rightMin := getMinChar(node.right)
		if rightMin < minChar {
			minChar = rightMin
		}
	}

	return minChar
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	item := x.(*Node)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return node
}

func (pq *PriorityQueue) Add(node *Node) {
	heap.Push(pq, node)
}

func (pq *PriorityQueue) Poll() *Node {
	return heap.Pop(pq).(*Node)
}

func buildPriorityQueue(frequencyMap *map[rune]int64) *PriorityQueue {
	// Create a slice to hold all entries
	type entry struct {
		char   rune
		weight int64
	}

	entries := make([]entry, 0, len(*frequencyMap))
	for char, weight := range *frequencyMap {
		entries = append(entries, entry{char, weight})
	}

	// Sort entries to ensure deterministic order
	// Sort by weight first, then by character
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].weight != entries[j].weight {
			return entries[i].weight < entries[j].weight
		}
		return entries[i].char < entries[j].char
	})

	// Build priority queue from sorted entries
	pq := make(PriorityQueue, len(entries))
	for i, e := range entries {
		pq[i] = newLeafNode(e.weight, e.char)
	}

	heap.Init(&pq)

	return &pq
}
