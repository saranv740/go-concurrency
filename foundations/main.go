package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var data int32
var l sync.Mutex

// Example of race condidtion, where multiple goroutines read same variable
// Note that memory access synchronization using mutex or atomic Addition prevents race condition
func RaceConditions() {
	for range 100 {
		go func() {
			l.Lock()
			data++
			fmt.Printf("%v\t", data)
			l.Unlock()
		}()
	}
	fmt.Printf("\n")
	time.Sleep(3 * time.Second)
	fmt.Println("done")
}

type Value struct {
	v  int
	mu sync.Mutex
}

// This function demonstrates the deadlock situation. where two go routines waiting for each other to unlock the variable.
// Running this will cause panic:
// fatal error: all goroutines are asleep - deadlock!
func DeadlockSimulation() {
	var wg sync.WaitGroup
	printSum := func(v1, v2 *Value) {
		defer wg.Done()
		v1.mu.Lock()
		defer v1.mu.Unlock()

		time.Sleep(2 * time.Second)

		v2.mu.Lock()
		defer v2.mu.Unlock()
		fmt.Println("sum ", v1.v+v2.v)
	}

	var a, b Value
	wg.Add(2)
	go printSum(&a, &b)
	go printSum(&b, &a)
	wg.Wait()
}

// This function demonstrates memory visiblity issue where loop might read inconsistent value from memory and
// it might go on forever. But when using atomic operations the memory is shared consistently
func MemoryVisibilityDemo() {
	var ready int32

	setReady := func() {
		fmt.Println("flg set")
		atomic.AddInt32(&ready, 1)
	}

	startLoop := func() {
		fmt.Println("loop")
		for atomic.LoadInt32(&ready) == 0 {

		}
		fmt.Println("loop end")
	}

	go startLoop()
	go setReady()

	time.Sleep(2 * time.Second)
	println("done")
}

func main() {
	// RaceConditions()
	// DeadlockSimulation()
	MemoryVisibilityDemo()
}
