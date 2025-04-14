package ratelimiting

import (
	"context"
	"log"
	"os"
	"sync"
)

//rate limiting can be of different types, it depends on various things
//in this case we try to implement the token bucket way, which has a given number of access token to start with
// and hence has a set burst value or the set of requests it can handle while warming up or at the  start
//next it has a certain limit defined after which the access token refreshed in the bucket for the next request to take it away
// so 2 things, d -> depth of the bucker , r-> the  rate at which the bucket is replenished

//simple client side apis

type APIConnection struct{}

func open() *APIConnection{
	return &APIConnection{}
}

func (a *APIConnection) ReadFile(ctx context.Context) error{

	//did some fake work here
	return nil
}

func ( a* APIConnection) ResolveAddress(ctx context.Context) error{
	//did some fake work here
	return nil
}



// 18:39:21 resolve address
// 18:39:21 resolve address
// 18:39:21 readfile
// 18:39:21 resolve address
// 18:39:21 readfile
// 18:39:21 resolve address
// 18:39:21 readfile
// 18:39:21 resolve address
// 18:39:21 readfile
// 18:39:21 resolve address
// 18:39:21 readfile
// 18:39:21 resolve address
// 18:39:21 readfile
// 18:39:21 readfile
// 18:39:21 readfile
// 18:39:21 resolve address
// 18:39:21 resolve address
// 18:39:21 readfile
// 18:39:21 readfile
// 18:39:21 resolve address
// 18:39:21 Done

// all of the requests are handled almost instantly, no rate limiter, and the requests are free to access the resources
func ClientSideNoRateLimiter(){

	defer log.Printf("Done")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConnection := open()
	var wg sync.WaitGroup
	wg.Add(20)

	for i:=0;i<10;i++{
		go func() {
			defer wg.Done()
			err := apiConnection.ReadFile(context.Background())
			if err != nil{
				log.Printf("Big reading error")
			}
			log.Printf("readfile")
		}()
	
	}
	
	for i:=0;i<10;i++{
		go func() {
			defer wg.Done()
			err := apiConnection.ReadFile(context.Background())
			if err != nil{
				log.Printf("Big resolving error")
			}
			log.Printf("resolve address")
		}()	
	
	}

	wg.Wait()
}

