package main

import (
	"fmt"
	"sync"
	"time"
)

type value struct {
	mu sync.Mutex
	val int
}

var wg sync.WaitGroup

func printsum(v1, v2 *value){
	defer wg.Done()

	v1.mu.Lock()
	defer v1.mu.Unlock()


	time.Sleep(2*time.Second)

	v2.mu.Lock()
	defer v2.mu.Unlock()

	fmt.Printf("sum=%v\n",v1.val + v2.val)
}
func main() {

	var a, b value
	wg.Add(2)
	go printsum(&a, &b)
	go printsum(&b, &a)
	wg.Wait()
}

// both the process try to acquire the lock for the v1 and v2 from each other and hence in a deadlock state