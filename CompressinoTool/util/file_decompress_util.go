package util

import (
	"compressiontool/huffmantree"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func DecompressFile(inputFile, outputFile string) error {
	file, err := os.ReadFile(inputFile)

	log.Println("File reading is done")

	if err != nil {
		return err
	}

	frequencyMap, binaryString, err := extractDataFromFileContent(file)

	log.Printf("Frquency map and binary string, %v\n", len(frequencyMap))

	if err != nil {
		return err
	}

	huffmantree := huffmantree.NewHuffManTree(frequencyMap)

	log.Println("Huffman tree building is done")

	fileContent := huffmantree.Decode(binaryString)

	err = writeDecodedContentToFile(outputFile, fileContent)

	log.Println("Wrote content to the file")

	if err != nil {
		return err
	}

	return nil
}

func unpackCompressedBytes(compressedBytes []byte, paddingLen int) string {
	var builder strings.Builder
	builder.Grow(len(compressedBytes) * 8)

	for _, b := range compressedBytes {
		builder.WriteString(fmt.Sprintf("%08b", b))
	}

	bitString := builder.String()

	if paddingLen > 0 {
		bitString = bitString[:len(bitString)-paddingLen]
	}

	return bitString
}

func extractDataFromFileContent(file []byte) (map[rune]int64, string, error) {
	headerSectionStartIndex := 0
	for i := range file {
		if file[i] == '\n' {
			headerSectionStartIndex = i
			break
		}
	}
	headerSectionLen, err := strconv.ParseInt(string(file[0:headerSectionStartIndex]), 10, 64)

	if err != nil {
		return nil, "", err
	}

	headerSectionStartIndex++

	paddingLen, err := strconv.ParseInt(string(file[headerSectionStartIndex:headerSectionStartIndex+1]), 10, 32)

	if err != nil {
		return nil, "", err
	}

	headerSectionStartIndex += 2

	headerSectionEndIndex := headerSectionStartIndex + int(headerSectionLen)

	headerSection := file[headerSectionStartIndex:headerSectionEndIndex]

	frequencyMap := make(map[rune]int64)

	err = json.Unmarshal(headerSection, &frequencyMap)

	if err != nil {
		return nil, "", err
	}

	compressedTextByteArr := file[headerSectionEndIndex:]

	return frequencyMap, unpackCompressedBytes(compressedTextByteArr, int(paddingLen)), nil
}

func writeDecodedContentToFile(filePath, content string) error {
	err := os.WriteFile(filePath, []byte(content), 0666)
	return err
}
