package heartbeat

import (
	"fmt"
	"time"
)

//heartbeat is to the signal that the node or the goroutine is not dead yet
//it is also a way to convey the meta information about the node
//there are two different types of hearbeat

//1. heartbeat that occur ona time interval
//2. heartbeat that occur at the bginning of the unit of work

//heartbeart is way for the the goroutine to signal that is waiting for a condition or another signal or a unit of work and is not dead yet

func DoWork(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time){
	heartbeat := make(chan interface{})
	results := make(chan time.Time)

	go func(){
		defer close(heartbeat)
		defer close(results)

		pulse := time.Tick(pulseInterval) // pulse for the hearbeat
		workGen := time.Tick(2*pulseInterval) // twice the duration to stimulate the work done with this interval

		sendPulse := func(){
			select{
			case heartbeat <-struct{}{}:
			default:
			}
		}

		// the select statement are within the for loops because, you might send multiple heartbeats before recieving any work, or send multiple work signals before recieving any heartbeat
		sendResult := func(r time.Time){
			for{
				select{
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case results<-r:
					return
				}
			}
		}


		for{
			select{
			case <-done:
				return
			case <-pulse:
				sendPulse()
			case r := <-workGen:
				sendResult(r)
			}
		}
	}()

	return heartbeat, results
}


//basic heartbeat send two pulses per second and then then a result pulse, 10 secodn the done channel cancels the flow

// Pulse
// Pulse
// results 42
// Pulse
// Pulse
// results 44
// Pulse
// Pulse
// results 46
// Pulse
// Pulse
// results 48
// Pulse
// Pulse

func BasicHeartbeart(){
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func(){close(done)})

	const timeout = 2*time.Second
	heartbeat, results := DoWork(done, timeout/2)
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
			return
		}
	}
}

