package sync

import (
	"fmt"
	"sync"
	"time"
)

//a  basic queue which has a size of 2, but we want to add 10 items only
//when the queue is empty
func BasicCond() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{},0 , 10)

	removeFromQueue := func(delay time.Duration){
		time .Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("Item Removed")
		c.L.Unlock()
		c.Signal()
	}

	for i:=0;i<10;i++{
		c.L.Lock()
		for len(queue) == 2{
			c.Wait()
		}
		fmt.Println("Adding to the queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(2*time.Second)
		c.L.Unlock()
	}
}

// internally the runtime maintains a FIFO list of the  goroutines waiting to be signalled
// and the go routine which has been waiting for the longest should be given the access
// Broadcast just informs all in the queue for the condition has been satisfied

// this is a simpla handler , for a button clicking which multiple functions should be triggered
func BroadcastCond(){

	type Button struct{
		Clicked *sync.Cond
	}

	button := Button{ Clicked : sync.NewCond(&sync.Mutex{})}
	

	subscribe :=  func(c * sync.Cond, fn func()){
		var goRoutineRunning sync.WaitGroup
		goRoutineRunning.Add(1)
		go func(){
			goRoutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()

		goRoutineRunning.Wait()
	}


	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)
	subscribe(button.Clicked, func() {
		fmt.Println("Maximizing Window")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying the dialog window")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Mouse Clicked")
		clickRegistered.Done()
	})

	button.Clicked.Broadcast()

	clickRegistered.Wait()
}
