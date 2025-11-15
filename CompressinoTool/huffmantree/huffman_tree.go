package huffmantree

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
		node1 := pq.Poll()
		node2 := pq.Poll()

		sum := node1.weight + node2.weight

		mergedNode := MergeTree(sum, node1, node2)

		pq.Add(mergedNode)
	}

	ht.root = pq.Poll()
}
