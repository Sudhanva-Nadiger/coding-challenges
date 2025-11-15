package huffmantree

import "container/heap"

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].weight < pq[j].weight
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

func BuildPriorityQueue(frequencyMap *map[rune]int64) *PriorityQueue {
	pq := make(PriorityQueue, len(*frequencyMap))

	i := 0

	for char, weight := range *frequencyMap {
		pq[i] = NewLeafNode(weight, char)
		i++
	}

	heap.Init(&pq)

	return &pq
}
