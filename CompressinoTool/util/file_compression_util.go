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

	frquencyMap := make(map[rune]int64)

	for _, char := range string(inputContent) {
		frquencyMap[char]++
	}

	huffmanTree := huffmantree.NewHuffManTree(frquencyMap)

	compressedBitString := huffmanTree.GetCompressedBitsString(string(inputContent))

	compressedByteArray, paddingLen := packCompressedBitString(compressedBitString)

	headerSection, err := getHeaderSectionToWrite(*huffmanTree.FrequencyMap)

	if err != nil {
		return err
	}

	err = writeToFile(headerSection, compressedByteArray, paddingLen, outputFile)

	logStats(inputFile, outputFile)

	if err != nil {
		return err
	}

	return nil
}

func packCompressedBitString(compressedBitString string) ([]byte, int) {
	paddedBitString := compressedBitString + strings.Repeat("0", 8-len(compressedBitString)%8)

	paddingBitsLen := len(paddedBitString) - len(compressedBitString)

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

	// write len of header section with new line
	err := os.WriteFile(outputFile, []byte(strconv.Itoa(len(headerSection))+"\n"), 0644)

	if err != nil {
		return err
	}

	// write padding length
	err = os.WriteFile(outputFile, []byte(strconv.Itoa(paddingLen)+"\n"), 0644)
	if err != nil {
		return err
	}

	// write header section
	err = os.WriteFile(outputFile, headerSection, 0644)
	if err != nil {
		return err
	}

	// write compressed byte array
	err = os.WriteFile(outputFile, compressedByteArray, 0644)
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
