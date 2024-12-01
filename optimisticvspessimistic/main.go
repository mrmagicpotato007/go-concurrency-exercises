package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var counter int64 = 1

type atomicCounter struct {
	counter int
	mu      sync.RWMutex
}

var ac *atomicCounter

func init() {
	ac = &atomicCounter{counter: 0}
}

func incrementCounter(goRoutineNumber int, wg *sync.WaitGroup) {
	// Decrements the counter when the goroutine completes
	defer wg.Done()

	oldCounterValue := atomic.LoadInt64(&counter)
	newCounterValue := oldCounterValue + 1 // perform the operation

	// CompareAndSwap operation is an atomic instruction
	// It only updates the counter = newCounterValue only if the current value is oldCounterValue
	result := atomic.CompareAndSwapInt64(&counter, oldCounterValue, newCounterValue)

	if !result {
		//fmt.Printf("%d operation failed\n", goRoutineNumber)
	} else {
		// fmt.Println("operation success")
	}
}

func main() {
	// Wait Group representing a pool of goroutines
	var wg sync.WaitGroup

	// for loop to simulate multiple threads updating the database at the same time
	startTime := time.Now()
	for i := 1; i <= 10000000; i++ {
		wg.Add(1)
		go incrementCounter(i, &wg)
	}

	wg.Wait()
	fmt.Printf("time taken for optimistic locking %s \n", time.Since(startTime).String())

	startTime = time.Now()
	for i := 1; i <= 10000000; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			ac.mu.Lock()
			defer ac.mu.Unlock()
			ac.counter++
		}(&wg)
	}
	wg.Wait()
	fmt.Printf("time taken for pessimistic locking %s \n", time.Since(startTime).String())
	fmt.Printf("atomic counter value %d \n", ac.counter)
}
