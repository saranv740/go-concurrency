package main

import (
	"fmt"
	"sync"
	"time"
)

func WaitGroupExample() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(salutation)
		}()
	}
	wg.Wait()
}

type Button struct {
	Clicked *sync.Cond
}

func CondExample() {
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func(c *sync.Cond, fn func()) {
		var isRunning sync.WaitGroup
		isRunning.Add(1)
		go func() {
			isRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		isRunning.Wait()
	}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)

	subscribe(button.Clicked, func() {
		fmt.Println("Maximizing window.")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Mouse clicked.")
		clickRegistered.Done()
	})

	button.Clicked.Broadcast()
	clickRegistered.Wait()
}

func ChannelsExample() {
	numberStream := make(chan int, 4)

	go func() {
		defer close(numberStream)
		for i := range 50 {
			numberStream <- (i + 1) * 2
		}
	}()

	// for i := range numberStream {
	// 	fmt.Printf("received number from the stream %d\n", i)
	// }

	v1, ok1 := <-numberStream
	v2, ok2 := <-numberStream
	v3, ok3 := <-numberStream

	fmt.Printf("val %d ok %v \n", v1, ok1)
	fmt.Printf("val %d ok %v \n", v2, ok2)
	fmt.Printf("val %d ok %v \n", v3, ok3)
}

func SelectExample() {
	start := time.Now()
	c := make(chan interface{})

	go func() {
		time.Sleep(5 * time.Second)
		close(c)
	}()

	fmt.Println("blocking on select")
	select {
	case <-c:
		fmt.Println("triggered after ", time.Since(start))
	}
}

func SelectWithMultipleChannels() {
	c1 := make(chan any)
	close(c1)
	c2 := make(chan any)
	close(c2)

	var c1Count, c2Count int

	for range 1000 {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}

	fmt.Println("c1 count ", c1Count, " c2 count ", c2Count)
}

func ForSelect() {
	done := make(chan any)
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	var workUnit int

loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}

		workUnit++
		time.Sleep(1 * time.Second)
	}

	fmt.Println("done ", workUnit, " unit of work")
}

func main() {
	WaitGroupExample()
	CondExample()
	ChannelsExample()
	SelectExample()
	SelectWithMultipleChannels()
	ForSelect()
}
