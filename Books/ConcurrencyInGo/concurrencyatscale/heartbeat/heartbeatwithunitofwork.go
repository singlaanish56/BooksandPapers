package heartbeat

import (
	"fmt"
	"math/rand"
)

func DoWorkWithUnitOfWork(done <-chan interface{}) (<-chan interface{}, <-chan int) {
	heartbeatStream := make(chan interface{}, 1) // always make a demo heartbeat of 1 pulse
	workStream := make(chan int)

	go func() {
		defer close(heartbeatStream)
		defer close(workStream)

		for i := 0; i < 10; i++ {
			select {
			case heartbeatStream <- struct{}{}:
			default:

			}

			select {
			case <-done:
				return
			case workStream <- rand.Intn(10):

			}
		}
	}()


	return heartbeatStream, workStream
}


func UnitWorkHeartbeat(){
	done := make(chan interface{})
	defer close(done)

	heartbeatStream, workStream := DoWorkWithUnitOfWork(done)
	for{
		select{
		case _,ok:= <-heartbeatStream:
			if ok{
				fmt.Println("Pulse")
			}else{
				return
			}
		case r, ok := <- workStream:
			if ok{
				fmt.Printf("results %v\n",r)
			}else{
				return
			}
		}
	}
}