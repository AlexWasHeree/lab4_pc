package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return nil, err
	}
	return data, nil
}

func chunkFile(data []byte, chunkSize int) []int {
	var chunks []int

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}

		chunk := data[i:end]
		chunkSum := 0
		for _, b := range chunk {
			chunkSum += int(b)
		}

		chunks = append(chunks, chunkSum)
	}

	return chunks
}

func compareChunks(chunks1, chunks2 []int) float64 {
	matchingChunks := 0
	totalChunks := len(chunks1)

	if len(chunks2) < totalChunks {
		totalChunks = len(chunks2)
	}

	for i := 0; i < totalChunks; i++ {
		if chunks1[i] == chunks2[i] {
			matchingChunks++
		}
	}

	return float64(matchingChunks) / float64(totalChunks)
}

func calculateSimilarities(filePaths []string, chunkSize int) {
	chunkedFiles := make(map[string][]int)

	for _, path := range filePaths {
		data, err := readFile(path)
		if err != nil {
			fmt.Printf("Skipping file %s due to error: %v\n", path, err)
			continue
		}

		chunks := chunkFile(data, chunkSize)
		chunkedFiles[path] = chunks
	}

	for i := 0; i < len(filePaths); i++ {
		file1 := filePaths[i]
		for j := i + 1; j < len(filePaths); j++ {
			file2 := filePaths[j]
			similarity := compareChunks(chunkedFiles[file1], chunkedFiles[file2])

			fmt.Printf("Similarity between %s and %s: %.6f%%\n", file1, file2, similarity*100)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	filePaths := os.Args[1:]
	chunkSize := 1024

	calculateSimilarities(filePaths, chunkSize)
}
