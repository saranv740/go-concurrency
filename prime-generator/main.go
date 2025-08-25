package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

const NUM_WORKERS = 4

func SieveOfEratosthenes(n int, results []bool) {
	limit := int(math.Sqrt(float64(n)))
	for i := 2; i <= limit; i++ {
		if results[i] == true {
			for j := i * i; j <= n; j += i {
				results[j] = false
			}
		}
	}
}

func BasePrimes(limit int) []int {
	baseArr := make([]bool, limit+1)
	for i := range baseArr {
		baseArr[i] = true
	}
	baseArr[0] = false
	baseArr[1] = false

	for i := 2; i*i <= limit; i++ {
		if baseArr[i] == true {
			for j := i * i; j <= limit; j += i {
				baseArr[j] = false
			}
		}
	}

	basePrimes := make([]int, 0)
	for idx, val := range baseArr {
		if val == true {
			basePrimes = append(basePrimes, idx)
		}
	}
	return basePrimes
}

func SieveOfEratosthenesConcurrent(limit int, results []bool) {
	root := int(math.Sqrt(float64(limit)))
	base := BasePrimes(root)

	workers := min(NUM_WORKERS, limit)
	var wg sync.WaitGroup

	chunkSize := limit / workers
	remainder := limit % workers

	for i := range NUM_WORKERS {
		low := (i * chunkSize) + 1
		high := (i + 1) * chunkSize

		if i == NUM_WORKERS-1 {
			high += remainder
		}

		wg.Add(1)
		go func(low, high int) {
			defer wg.Done()
			for _, prime := range base {
				// if prime = 7 and low = 100 we need first multiple of 7 above 100 so 100/7 = 14 and 14*7 = 98. 98+7 = 105
				start := (low / prime) * prime
				if start < low {
					start += prime
				}

				// if low is 20 and we are getting prime as 7 we don't need to start from 21 those will be marked as false
				// by other primes such as 3
				if start < prime*prime {
					start = prime * prime
				}

				for j := start; j <= high; j += prime {
					results[j] = false
				}
			}
		}(low, high)
	}
	wg.Wait()
}

func main() {
	n := 25_00_003

	// Sequential implementation
	results := make([]bool, n+1)
	for i := range results {
		results[i] = true
	}

	start := time.Now()

	results[0] = false
	results[1] = false
	SieveOfEratosthenes(n, results)

	fmt.Printf("time taken for n=%d is %v\n", n, time.Since(start))

	// concurrernt implementation
	resultCon := make([]bool, n+1)
	for i := range resultCon {
		resultCon[i] = true
	}

	start = time.Now()

	resultCon[0] = false
	resultCon[1] = false
	SieveOfEratosthenesConcurrent(n, resultCon)

	fmt.Printf("time taken for concurrent processing n=%d is %v\n", n, time.Since(start))

	for idx := range results {
		if results[idx] != resultCon[idx] {
			fmt.Printf("Mismatch at index %d, result[%d]=%v, resultConcurrent[%d]=%v", idx, idx, results[idx], idx, resultCon[idx])
		}
	}
}
