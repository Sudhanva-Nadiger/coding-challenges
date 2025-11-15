package util

import (
	"compressiontool/huffmantree"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
)

func CreateCompressedFile(inputFile string, outputFile string) error {
	inputContent, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	log.Println("reading input file compreted")

	frquencyMap := make(map[rune]int64)

	for _, char := range string(inputContent) {
		frquencyMap[char]++
	}

	log.Printf("frequency map creation is done, %v\n", len(frquencyMap))

	huffmanTree := huffmantree.NewHuffManTree(frquencyMap)

	log.Println("Huffman tree building is done")

	compressedBitString := huffmanTree.Encode(string(inputContent))

	log.Println("Content encoding is done")

	compressedByteArray, paddingLen := packCompressedBitString(compressedBitString)

	log.Println("packing is done")

	headerSection, err := getHeaderSectionToWrite(*huffmanTree.FrequencyMap)

	if err != nil {
		return err
	}

	err = writeToFile(headerSection, compressedByteArray, paddingLen, outputFile)

	log.Println("Wrote to file")

	logStats(inputFile, outputFile)

	if err != nil {
		return err
	}

	return nil
}

func packCompressedBitString(compressedBitString string) ([]byte, int) {
	// Calculate padding needed (0 if already multiple of 8)
	paddingBitsLen := (8 - len(compressedBitString)%8) % 8

	paddedBitString := compressedBitString + strings.Repeat("0", paddingBitsLen)

	packedBytes := make([]byte, len(paddedBitString)/8)

	for i := range packedBytes {
		byteValue := paddedBitString[i*8 : (i+1)*8]
		byteValueInt, _ := strconv.ParseUint(byteValue, 2, 8)
		packedBytes[i] = byte(byteValueInt)
	}

	return packedBytes, paddingBitsLen
}

func getHeaderSectionToWrite(frequencyMap map[rune]int64) ([]byte, error) {

	jsonBytes, err := json.Marshal(frequencyMap)

	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}

func writeToFile(headerSection, compressedByteArray []byte, paddingLen int, outputFile string) error {
	headerLenStr := strconv.Itoa(len(headerSection)) + "\n"
	paddingLenStr := strconv.Itoa(paddingLen) + "\n"

	totalSize := len(headerLenStr) + len(paddingLenStr) + len(headerSection) + len(compressedByteArray)

	fileContent := make([]byte, 0, totalSize)
	fileContent = append(fileContent, []byte(headerLenStr)...)
	fileContent = append(fileContent, []byte(paddingLenStr)...)
	fileContent = append(fileContent, headerSection...)
	fileContent = append(fileContent, compressedByteArray...)

	// Write once
	err := os.WriteFile(outputFile, fileContent, 0644)
	if err != nil {
		return err
	}

	return nil
}

func logStats(inputFileName, outputFileName string) {
	inputFile, _ := os.Open(inputFileName)
	outputFile, _ := os.Open(outputFileName)

	defer inputFile.Close()
	defer outputFile.Close()

	inputFileStat, _ := inputFile.Stat()
	outputFileStat, _ := outputFile.Stat()

	log.Println("Input file size: ", inputFileStat.Size())
	log.Println("Compressed file size: ", outputFileStat.Size())

	compressionRatio := float64(inputFileStat.Size()) / float64(outputFileStat.Size())
	log.Printf("Compression ratio: %.2f%%\n", compressionRatio*100)
}
