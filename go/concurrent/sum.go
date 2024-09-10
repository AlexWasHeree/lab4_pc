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


func sum(filePath string) (int, error) {
	data, err := readFile(filePath)
	if err != nil {
		return 0, err
	}

	_sum := 0
	for _, b := range data {
		_sum += int(b)
	}

	return _sum, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	var wg sync.WaitGroup
	sumChannel := make(chan struct {
		sum  int
		path string
	}, len(os.Args[1:])) 


	sums := make(map[int][]string)
	var totalSum int64
	var mu sync.Mutex /

	for _, path := range os.Args[1:] {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()

			_sum, err := sum(filePath)
			if err != nil {
				return
			}

			sumChannel <- struct {
				sum  int
				path string
			}{_sum, filePath}
		}(path)
	}

	go func() {
		wg.Wait()
		close(sumChannel)
	}()

	for result := range sumChannel {
		mu.Lock()
		totalSum += int64(result.sum)
		sums[result.sum] = append(sums[result.sum], result.path)
		mu.Unlock()
	}

	fmt.Println("Concurrent: ", totalSum)

	for sum, files := range sums {
		if len(files) > 1 {
			fmt.Printf("Sum %d: %v\n", sum, files)
		}
	}
}
