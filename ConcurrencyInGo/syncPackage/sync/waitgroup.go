package sync

import (
	"fmt"
	"sync"
	"time"
)

//helps with low level memory sync access

//wait group
//wait for the set of the concurrent actions to complete

func WaitGroup(){
	var wg sync.WaitGroup

	wg.Add(1)
	go func(){
		defer wg.Done()
		fmt.Println("yooohoooo sleep time")
		time.Sleep(2*time.Second)
	}()

	wg.Wait()
	fmt.Println("all complete")
}
