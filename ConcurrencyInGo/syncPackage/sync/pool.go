package sync

import (
	"fmt"
	"sync"
)

func SimplePool() {
	myPool := &sync.Pool{
		New : func() interface{}{
			fmt.Println("Creating new instance")
			return struct{}{}
		},
	}

	myPool.Get()
	instance := myPool.Get()
	myPool.Put(instance)
	myPool.Get()
}