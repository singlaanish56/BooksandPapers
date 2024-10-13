package sync

import (
	"fmt"
	"sync"
)

//what once does is tracks the number of the calls to he DO function
//doesnt matter the function passed inside the DO function
//

func Onceonce(){
	var count int
	increment := func(){
		count++
	}
	
	
	
	var once sync.Once

	var increments sync.WaitGroup
	increments.Add(100)
	for i:=0;i<100;i++{
		go func(){
			defer increments.Done()
			once.Do(increment)
		}()
	}

	increments.Wait()
	fmt.Printf("The number of times is called %d\n", count)
}