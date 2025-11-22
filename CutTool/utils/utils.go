package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var builder strings.Builder

func ParseFields(s string) ([]int64, error) {
	fields := []int64{}

	splitString := ""

	s = strings.Trim(s, " ,")

	if strings.Contains(s, ",") {
		splitString = ","
	} else if strings.Contains(s, " ") {
		splitString = " "
	}

	nums := strings.SplitSeq(s, splitString)

	for num := range nums {
		index, err := strconv.ParseInt(num, 10, 64)
		if err != nil {
			return nil, err
		}
		fields = append(fields, index)
	}

	return fields, nil
}

func ProcessLine(fields *[]int64, line, delimiter *string) string {
	words := strings.Split(*line, *delimiter)
	fieldsLen := len(*fields)

	builder.Grow(len(*line))
	defer builder.Reset()

	for i, index := range *fields {
		builder.WriteString(words[index])
		if i != fieldsLen-1 {
			builder.WriteString(*delimiter)
		}
	}

	return builder.String()
}

func GetScanner(fileName string) (*bufio.Scanner, error) {
	var scanner *bufio.Scanner

	if fileName == "" || fileName == "-" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		file, err := os.Open(fileName)

		if err != nil {
			return nil, fmt.Errorf("error reading file %v", err)
		}

		defer file.Close()

		scanner = bufio.NewScanner(file)
	}

	return scanner, nil
}
