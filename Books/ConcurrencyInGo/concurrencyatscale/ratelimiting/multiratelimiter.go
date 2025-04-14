package ratelimiting

import (
	"context"
	"sort"
	"time"
	"log"
	"os"
	"sync"
	"golang.org/x/time/rate"
)


//measures it in terms of the operations per time , rather the time in between the events
func Per(eventcount int, duration time.Duration) rate.Limit{
	return rate.Every(duration/time.Duration(eventcount))
}

//now there could be multu rate limiting, per second, per minute, per hour
// these multi tiers, fine granular to limit requests per second, to coarse grained request

// so the we can define the rate limit recursively
type RateLimiter interface{
	Wait(context.Context) error
	Limit() rate.Limit
}

type multimeter struct{
	limiters []RateLimiter
}

func MultiLimiter(limiters... RateLimiter) *multimeter{
	bylimt := func(i, j int) bool{
		return limiters[i].Limit() < limiters[j].Limit()
	}

	sort.Slice(limiters, bylimt)
	return &multimeter{limiters: limiters}
}

func (l *multimeter) Wait(ctx context.Context) error{
	for _, l := range l.limiters{
		if err := l.Wait(ctx); err!=nil{
			return err
		}
	}

	return nil
}

//returns the most restrictive limit , which  is the first one , because we sort it while creating the multiLimiter
func (l * multimeter) Limit() rate.Limit{
	return l.limiters[0].Limit()
}


type apiConnectionMultiRateLimit struct{
	ratelimiter RateLimiter
}

//the second is 2 request per second, ith no burstiness, the minutes is 10 request per minute and with a burstiness of 10
func open3() *apiConnectionMultiRateLimit{
	secondsLimit := rate.NewLimiter(Per(2, time.Second), 2)
	minutesLimit := rate.NewLimiter(Per(10, time.Minute), 10)
	return &apiConnectionMultiRateLimit{
		ratelimiter: MultiLimiter(secondsLimit, minutesLimit),
	}
}

func (a * apiConnectionMultiRateLimit) ReadFile3(ctx context.Context) error{
	if err := a.ratelimiter.Wait(ctx); err != nil{
		return err
	}

	return nil
}

func (a *apiConnectionMultiRateLimit) ResolveAddress3(ctx context.Context) error{
	if err := a.ratelimiter.Wait(ctx); err != nil{
		return err
	}

	return nil
}


// 20:26:22 read file
// 20:26:22 resolve address
// 20:26:23 read file
// 20:26:23 read file
// 20:26:24 read file
// 20:26:24 read file
// 20:26:25 resolve address
// 20:26:25 resolve address
// 20:26:26 resolve address
// 20:26:26 resolve address
// 20:26:28 resolve address
// 20:26:34 resolve address
// 20:26:40 resolve address
// 20:26:46 resolve address
// 20:26:52 resolve address
// 20:26:58 read file
// 20:27:04 read file
// 20:27:10 read file
// 20:27:16 read file
// 20:27:22 read file
// 20:27:22 done

// 11 request for 2 requests per second, and then the token refreshed every 6 second 10/60 seconds, 1 per 6 seconds
func ClientSideMultiLimiter(){
	defer log.Printf("done")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiconnection := open3()
	var wg sync.WaitGroup
	wg.Add(20)

	for i:=0;i<10;i++{
		go func() {
			defer wg.Done()
			err := apiconnection.ReadFile3(context.Background())
			if err != nil{
				log.Printf("err from the read file %v\n", err)
			}

			log.Printf("read file")
		}()
	}

	for i:=0;i<10;i++{
		go func() {
			defer wg.Done()
			err := apiconnection.ResolveAddress3(context.Background())
			if err !=nil{
				log.Printf("err from the resolve %v\n", err)
			}
			log.Printf("resolve address")
		} ()
	}



	wg.Wait()	
}