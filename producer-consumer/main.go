package main

import (
	"fmt"
	"sync"
)

func main() {
	N := 10
	stringStream := make(chan string)

	var wg sync.WaitGroup

	multiplesProducer := func(number, limit int) {
		defer wg.Done()
		for i := 1; i <= limit; i++ {
			stringStream <- fmt.Sprintf("mutiple %d * %d = %d", i, number, i*number)
		}
	}

	for i := 1; i <= N; i++ {
		wg.Add(1)
		go multiplesProducer(i, 10)
	}

	go func() {
		wg.Wait()
		close(stringStream)
	}()

	for val := range stringStream {
		fmt.Println(val)
	}

	fmt.Println("program done")
}
