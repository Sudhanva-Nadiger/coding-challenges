package huffmantree

import (
	"container/heap"
	"testing"
)

func TestPriorityQueue_Len(t *testing.T) {
	pq := PriorityQueue{
		newLeafNode(5, 'a'),
		newLeafNode(3, 'b'),
	}

	if pq.Len() != 2 {
		t.Errorf("Expected length 2, got %d", pq.Len())
	}
}

func TestPriorityQueue_Less_DifferentWeights(t *testing.T) {
	pq := PriorityQueue{
		newLeafNode(5, 'a'),
		newLeafNode(3, 'b'),
	}

	// Node with weight 3 should be less than node with weight 5
	if !pq.Less(1, 0) {
		t.Error("Expected node with weight 3 to be less than node with weight 5")
	}
}

func TestPriorityQueue_Less_SameWeight_LeafNodes(t *testing.T) {
	pq := PriorityQueue{
		newLeafNode(5, 'b'),
		newLeafNode(5, 'a'),
	}

	// With same weight, 'a' should be less than 'b'
	if !pq.Less(1, 0) {
		t.Error("Expected node with char 'a' to be less than node with char 'b' when weights are equal")
	}
}

func TestPriorityQueue_Less_SameWeight_LeafVsInternal(t *testing.T) {
	leaf := newLeafNode(5, 'a')
	internal := &Node{
		weight: 5,
		left:   newLeafNode(2, 'b'),
		right:  newLeafNode(3, 'c'),
	}

	pq := PriorityQueue{internal, leaf}

	// Leaf node should come before internal node with same weight
	if !pq.Less(1, 0) {
		t.Error("Expected leaf node to come before internal node with same weight")
	}
}

func TestPriorityQueue_Swap(t *testing.T) {
	node1 := newLeafNode(5, 'a')
	node2 := newLeafNode(3, 'b')
	pq := PriorityQueue{node1, node2}

	pq.Swap(0, 1)

	if pq[0] != node2 || pq[1] != node1 {
		t.Error("Swap did not work correctly")
	}
}

func TestPriorityQueue_PushPop(t *testing.T) {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	node1 := newLeafNode(5, 'a')
	node2 := newLeafNode(3, 'b')
	node3 := newLeafNode(7, 'c')

	heap.Push(&pq, node1)
	heap.Push(&pq, node2)
	heap.Push(&pq, node3)

	if pq.Len() != 3 {
		t.Errorf("Expected length 3, got %d", pq.Len())
	}

	// Should pop in order: 3, 5, 7
	first := heap.Pop(&pq).(*Node)
	if first.weight != 3 || first.char != 'b' {
		t.Errorf("Expected first popped node to have weight 3 and char 'b', got weight %d and char %c", first.weight, first.char)
	}

	second := heap.Pop(&pq).(*Node)
	if second.weight != 5 || second.char != 'a' {
		t.Errorf("Expected second popped node to have weight 5 and char 'a', got weight %d and char %c", second.weight, second.char)
	}

	third := heap.Pop(&pq).(*Node)
	if third.weight != 7 || third.char != 'c' {
		t.Errorf("Expected third popped node to have weight 7 and char 'c', got weight %d and char %c", third.weight, third.char)
	}
}

func TestPriorityQueue_AddPoll(t *testing.T) {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	node1 := newLeafNode(5, 'a')
	node2 := newLeafNode(3, 'b')

	pq.Add(node1)
	pq.Add(node2)

	first := pq.Poll()
	if first.weight != 3 {
		t.Errorf("Expected first polled weight 3, got %d", first.weight)
	}

	second := pq.Poll()
	if second.weight != 5 {
		t.Errorf("Expected second polled weight 5, got %d", second.weight)
	}
}

func TestGetMinChar_LeafNode(t *testing.T) {
	node := newLeafNode(5, 'x')
	minChar := getMinChar(node)

	if minChar != 'x' {
		t.Errorf("Expected minChar 'x', got %c", minChar)
	}
}

func TestGetMinChar_InternalNode(t *testing.T) {
	// Tree structure:
	//       *
	//      / \
	//     c   *
	//        / \
	//       a   z
	node := &Node{
		weight: 10,
		left:   newLeafNode(3, 'c'),
		right: &Node{
			weight: 7,
			left:   newLeafNode(2, 'a'),
			right:  newLeafNode(5, 'z'),
		},
	}

	minChar := getMinChar(node)

	if minChar != 'a' {
		t.Errorf("Expected minChar 'a', got %c", minChar)
	}
}

func TestBuildPriorityQueue_Deterministic(t *testing.T) {
	frequencyMap := map[rune]int64{
		'a': 5,
		'b': 3,
		'c': 3,
		'd': 1,
	}

	// Build the priority queue multiple times
	pq1 := buildPriorityQueue(&frequencyMap)
	pq2 := buildPriorityQueue(&frequencyMap)

	// Both queues should have the same order
	for i := 0; i < pq1.Len(); i++ {
		node1 := heap.Pop(pq1).(*Node)
		node2 := heap.Pop(pq2).(*Node)

		if node1.weight != node2.weight || node1.char != node2.char {
			t.Errorf("Priority queues are not deterministic at index %d: got (%d, %c) and (%d, %c)",
				i, node1.weight, node1.char, node2.weight, node2.char)
		}
	}
}

func TestBuildPriorityQueue_SortedOrder(t *testing.T) {
	frequencyMap := map[rune]int64{
		'd': 1,
		'b': 3,
		'c': 3,
		'a': 5,
	}

	pq := buildPriorityQueue(&frequencyMap)

	// Expected order after heap operations:
	// 1. d(1) - lowest weight
	// 2. b(3) - same weight as c, but 'b' < 'c'
	// 3. c(3)
	// 4. a(5) - highest weight

	expected := []struct {
		weight int64
		char   rune
	}{
		{1, 'd'},
		{3, 'b'},
		{3, 'c'},
		{5, 'a'},
	}

	for i, exp := range expected {
		node := heap.Pop(pq).(*Node)
		if node.weight != exp.weight || node.char != exp.char {
			t.Errorf("At position %d: expected (%d, %c), got (%d, %c)",
				i, exp.weight, exp.char, node.weight, node.char)
		}
	}
}

func TestPriorityQueue_DeterministicWithMergedNodes(t *testing.T) {
	// Test that merged nodes also maintain deterministic order
	frequencyMap := map[rune]int64{
		'a': 1,
		'b': 1,
		'c': 2,
		'd': 2,
	}

	pq1 := buildPriorityQueue(&frequencyMap)
	pq2 := buildPriorityQueue(&frequencyMap)

	// Simulate Huffman tree building
	for i := 0; i < 2; i++ {
		// Build tree 1
		node1_1 := pq1.Poll()
		node1_2 := pq1.Poll()
		merged1 := &Node{
			weight: node1_1.weight + node1_2.weight,
			left:   node1_1,
			right:  node1_2,
		}
		pq1.Add(merged1)

		// Build tree 2
		node2_1 := pq2.Poll()
		node2_2 := pq2.Poll()
		merged2 := &Node{
			weight: node2_1.weight + node2_2.weight,
			left:   node2_1,
			right:  node2_2,
		}
		pq2.Add(merged2)

		// Check that both trees have same structure
		if merged1.weight != merged2.weight {
			t.Errorf("Iteration %d: merged nodes have different weights", i)
		}

		if merged1.left.char != merged2.left.char || merged1.right.char != merged2.right.char {
			t.Errorf("Iteration %d: merged nodes have different structure", i)
		}
	}
}
func TestPriorityQueue_UnicodeCharacters(t *testing.T) {
	frequencyMap := map[rune]int64{
		'é': 3,
		'ñ': 3,
		'ü': 1,
		'中': 5,
	}

	pq := buildPriorityQueue(&frequencyMap)

	// Should be sorted by weight first, then by unicode value
	first := pq.Poll()
	if first.weight != 1 {
		t.Errorf("Expected first weight 1, got %d", first.weight)
	}

	second := pq.Poll()
	third := pq.Poll()

	// Both have weight 3, should be sorted by char
	if second.weight != 3 || third.weight != 3 {
		t.Error("Expected both second and third to have weight 3")
	}

	if second.char >= third.char {
		t.Errorf("Expected chars to be sorted: got %c before %c", second.char, third.char)
	}
}
