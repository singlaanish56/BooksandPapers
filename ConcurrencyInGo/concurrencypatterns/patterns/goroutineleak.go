package patterns

import (
	"fmt"
	"time"
)

//leaking goroutine
//runtime doesnt kill the goroutines on its own

func LeakingRoutine(){
	doWork := func(strings <-chan string) <-chan interface{}{
		completed := make(chan interface{})

		go func(){
			defer fmt.Println("do work exited.")
			defer close(completed)

			for s:= range strings{
				fmt.Println(s)
			}
		}()

		return completed
	}

	doWork(nil)
	fmt.Println("Done .")
}

//how to fix this, we pass a done channel , which establishes the signal betweent the parent and the children subsroutine
// and hence allows for the cancellation of the go routine

func PreventLeakGoroutine(){
	doWork := func(done <-chan interface{},strings <-chan string) <- chan interface{}{
		terminated := make(chan interface{})
		go func(){
			defer fmt.Println("dowork exited")
			defer close(terminated)
			for {
				select{
				case s := <-strings:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()	

		return terminated
	}

	done := make(chan interface{})
	terminated := doWork(done,nil)

	go func(){
		time.Sleep(1*time.Second)
		fmt.Println("Cancelling the dowork goroutine")
		close(done)
	}()

	<-terminated
	fmt.Println("Done .")
}