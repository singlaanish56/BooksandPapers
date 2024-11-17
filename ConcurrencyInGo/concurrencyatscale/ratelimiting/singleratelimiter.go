package ratelimiting

import (
	"context"
	"log"
	"os"
	"sync"

	"golang.org/x/time/rate"
)

//this is a single rate limiter for the client side, obv rate limit can be done more permanently on the server side as well
//but this works when you dont want the client to send any requests either.
type apiConnectionRateLimit struct{
	ratelimiter *rate.Limiter
}

//1 event per second, burst of 1 token
func open2() *apiConnectionRateLimit{
	return &apiConnectionRateLimit{
		ratelimiter: rate.NewLimiter(rate.Limit(1), 1),
	}
}


//wait is a shorthand to Waitn, which returns an error if the burst size is exceeded, context is cancelled or deadline is missed
func (a *apiConnectionRateLimit) readFile2(ctx context.Context) error{
	if err := a.ratelimiter.Wait(ctx); err!=nil{
		return err
	}

	//pretend to do some work here

	return nil
}

func (a * apiConnectionRateLimit) resolveAddress2(ctx context.Context) error{
	if err := a.ratelimiter.Wait(ctx); err!=nil{
		return err
	}

	//pretend to do some work here

	return nil

}

// 19:53:45 read file
// 19:53:46 read file
// 19:53:47 resolve address
// 19:53:48 resolve address
// 19:53:49 resolve address
// 19:53:50 resolve address
// 19:53:51 resolve address
// 19:53:52 resolve address
// 19:53:53 resolve address
// 19:53:54 resolve address
// 19:53:55 read file
// 19:53:56 read file
// 19:53:57 read file
// 19:53:58 read file
// 19:53:59 read file
// 19:54:00 read file
// 19:54:01 read file
// 19:54:02 read file
// 19:54:03 resolve address
// 19:54:04 resolve address
// 19:54:04 done

// 1 request per second is served
func ClientSideSingleRateLimiter(){
	defer log.Printf("done")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiconnection := open2()
	var wg sync.WaitGroup
	wg.Add(20)

	for i:=0;i<10;i++{
		go func() {
			defer wg.Done()
			err := apiconnection.readFile2(context.Background())
			if err != nil{
				log.Printf("err from the read file %v\n", err)
			}

			log.Printf("read file")
		}()
	}

	for i:=0;i<10;i++{
		go func() {
			defer wg.Done()
			err := apiconnection.resolveAddress2(context.Background())
			if err !=nil{
				log.Printf("err from the resolve %v\n", err)
			}
			log.Printf("resolve address")
		} ()
	}



	wg.Wait()
}

