package heartbeat

import (
	"time"
	"fmt"
)

//this do work
//here we stop the goroutine after two tries to stimualte how the heatbeat helps us to identify a problem and close the process early
func HandicappedDoWork(done <- chan interface{}, pulseInterval time.Duration) (<- chan interface{}, <- chan time.Time){

	heartbeat := make(chan interface{})
	results := make(chan time.Time)

	go func(){
		defer close(heartbeat)
		defer close(results)
		pulse := time.Tick(pulseInterval)
		workGen := time.Tick(2*pulseInterval)

		sendPulse := func(){
			select{
			case heartbeat<-struct{}{}:
			default:
			}
		}

		sendResult := func(r time.Time){
			for{
				select{
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case results <- r:
					return
				}
			}
		}


		for i:=0;i<2;i++{
			select{
			case <-done:
				return
			case <-pulse:
				sendPulse()
			case r:= <-workGen:
				sendResult(r)
			}
		}
	}()

	return heartbeat, results
}

//so here the system realises that the there is issue with the hearbeat within the 2 second ans then times outs appropriately because nothing was recieving after that
func HeartbeartThatBreaks(){
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func(){close(done)})

	const timeout = 2*time.Second
	heartbeat, results := HandicappedDoWork(done, timeout/2)
	for{
		select{
		case _, ok:= <-heartbeat:
			if ok==false{
				return
			}
			fmt.Println("Pulse")
		case r, ok := <-results:
			if ok==false{
				return
			}
			fmt.Printf("results %v\n",r.Second())
		case <-time.After(timeout): //timeput if we dont recieve anew heartbeat or a result
			fmt.Println("woopsie issue witht the worker goroutine")	
			return
		}
	}
}

