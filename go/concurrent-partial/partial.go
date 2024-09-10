package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
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
	var mu sync.Mutex

	var wg sync.WaitGroup

	results := make(chan string, len(filePaths)*(len(filePaths)-1)/2)

	for _, path := range filePaths {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			data, err := readFile(filePath)
			if err != nil {
				fmt.Printf("Skipping file %s due to error: %v\n", filePath, err)
				return
			}

			chunks := chunkFile(data, chunkSize)

			mu.Lock()
			chunkedFiles[filePath] = chunks
			mu.Unlock()
		}(path)
	}

	wg.Wait()

	for i := 0; i < len(filePaths); i++ {
		file1 := filePaths[i]
		for j := i + 1; j < len(filePaths); j++ {
			file2 := filePaths[j]

			wg.Add(1)
			go func(file1, file2 string) {
				defer wg.Done()

				mu.Lock()
				chunks1 := chunkedFiles[file1]
				chunks2 := chunkedFiles[file2]
				mu.Unlock()

				
				similarity := compareChunks(chunks1, chunks2)

				
				results <- fmt.Sprintf("Similarity between %s and %s: %.6f%%", file1, file2, similarity*100)
			}(file1, file2)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
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
