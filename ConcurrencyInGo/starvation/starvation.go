package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup
var sharedLock sync.Mutex
const runtime = 1*time.Second

func greedyWorker(){
	defer wg.Done()

	var count int
	for begin := time.Now();time.Since(begin)<=runtime;{
		sharedLock.Lock()
		time.Sleep(3*time.Nanosecond)
		sharedLock.Unlock()
		count++
	}

	fmt.Printf("the greedy worker was able to the work in %v loops\n", count)
}

func politeworker(){
	defer wg.Done()

	var count int
	for begin := time.Now();time.Since(begin)<=runtime;{
		sharedLock.Lock()
		time.Sleep(1*time.Nanosecond)
		sharedLock.Unlock()

		sharedLock.Lock()
		time.Sleep(1*time.Nanosecond)
		sharedLock.Unlock()

		sharedLock.Lock()
		time.Sleep(1*time.Nanosecond)
		sharedLock.Unlock()

		count++
	}
	fmt.Printf("the polite worker was able to the work in %v loops\n", count)
}

func main() {
	wg.Add(2)


	go greedyWorker()
	go politeworker()
	wg.Wait()
}
