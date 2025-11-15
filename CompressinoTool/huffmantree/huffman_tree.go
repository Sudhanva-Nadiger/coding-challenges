package huffmantree

import (
	"fmt"
	"strings"
)

type HuffManTree struct {
	root            *Node
	prefixCodeTable *map[rune]string
	FrequencyMap    *map[rune]int64
}

func NewHuffManTree(frequencyMap map[rune]int64) *HuffManTree {
	ht := &HuffManTree{
		root:         nil,
		FrequencyMap: &frequencyMap,
	}

	ht.buildHuffManTree()

	return ht
}

func (ht *HuffManTree) Encode(inputContent string) string {
	ht.buildprefixCodeTable()

	var builder strings.Builder
	builder.Grow(len(inputContent) * 8)

	for _, char := range inputContent {
		builder.WriteString((*ht.prefixCodeTable)[char])
	}

	return builder.String()
}

func decodeBinaryString(binaryString *string, root *Node, index *int) rune {
	if root.isLeadNode() {
		return root.char
	}

	if (*binaryString)[*index] == '0' {
		*index += 1
		return decodeBinaryString(binaryString, root.left, index)
	}

	*index += 1

	return decodeBinaryString(binaryString, root.right, index)
}

func (ht *HuffManTree) Decode(binaryString string) string {
	index := 0
	len := len(binaryString)

	var builder strings.Builder
	builder.Grow(len / 8)

	for index < len {
		char := decodeBinaryString(&binaryString, ht.root, &index)
		builder.WriteRune(char)
	}

	return builder.String()
}

func (ht *HuffManTree) buildHuffManTree() {
	pq := buildPriorityQueue(ht.FrequencyMap)

	for pq.Len() > 1 {
		node1 := pq.Poll()
		node2 := pq.Poll()

		sum := node1.weight + node2.weight

		mergedNode := mergeTree(sum, node1, node2)

		pq.Add(mergedNode)
	}

	ht.root = pq.Poll()
}

func buildprefixCodeTableUtil(root *Node, prefixCode string, prefixCodeTable *map[rune]string) {
	if root == nil {
		return
	}

	if root.isLeadNode() {
		(*prefixCodeTable)[root.char] = prefixCode
		return
	}

	buildprefixCodeTableUtil(root.left, prefixCode+"0", prefixCodeTable)
	buildprefixCodeTableUtil(root.right, prefixCode+"1", prefixCodeTable)
}

func (ht *HuffManTree) buildprefixCodeTable() {
	prefixCodeTable := make(map[rune]string)
	ht.prefixCodeTable = &prefixCodeTable

	if ht.root == nil {
		return
	}

	// Start building the prefix code table from the root with empty prefix
	buildprefixCodeTableUtil(ht.root, "", &prefixCodeTable)
}

func (ht HuffManTree) PrintTree() {
	ht.root.PrintTree()

	fmt.Println()
}

func (ht HuffManTree) PrintprefixCodeTable() {
	if ht.prefixCodeTable == nil || len(*ht.prefixCodeTable) == 0 {
		fmt.Println("Prefix code table is empty")
		return
	}

	// Print table header
	fmt.Println("┌────────────┬──────────────┐")
	fmt.Println("│ Character  │ Huffman Code │")
	fmt.Println("├────────────┼──────────────┤")

	// Print each character and its code
	for char, code := range *ht.prefixCodeTable {
		// Format character for display
		var charDisplay string
		switch char {
		case '\n':
			charDisplay = "\\n"
		case '\t':
			charDisplay = "\\t"
		case '\r':
			charDisplay = "\\r"
		case ' ':
			charDisplay = "SPACE"
		default:
			charDisplay = string(char)
		}

		// Print row with proper padding
		fmt.Printf("│ %-10s │ %-12s │\n", charDisplay, code)
	}

	// Print table footer
	fmt.Println("└────────────┴──────────────┘")
}
