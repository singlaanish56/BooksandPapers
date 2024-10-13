package sync

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

//mutex
// mutual exclusion access to the critical section of the program
// shared resources can be accessed concurrently

func  Mutex(){
	 var count int
	 var lock sync.Mutex

	 increment := func(){
		lock.Lock()
		defer lock.Unlock()
		count++;
		fmt.Printf("in the increment call the value is %d\n", count)
	 }

	 decrement := func(){
		lock.Lock()
		defer lock.Unlock()
		count--;
		fmt.Printf("in the decrement call the value is %d\n", count)
	 }

	 var wg sync.WaitGroup
	 for i:=0;i<5;i++{
		wg.Add(1)
		go func(){
			defer wg.Done()
			increment()
		}()
	}

	for i:=0;i<5;i++{
		wg.Add(1)
		go func(){
			defer wg.Done()
			decrement()
		}()
	 }

	 wg.Wait()
	 fmt.Println("All complete")
}

//rw mutex
// now to define the access to the shared resources more granular we have some locks
// like the read / write lock (self explainable)
// this code has a slow producer and ever increasing count of observerm testing how much times it take for both the types of locks

// Readers  RWMutex    Mutex
// 1        77.1173ms  77.3886ms
// 2  77.5878ms  76.8515ms
// 4  79.7938ms  79.2016ms
// 8  78.3069ms  77.3034ms
// 16  77.1637ms  77.477ms
// 32  76.3619ms  77.2942ms
// 64  76.693ms  76.0862ms
// 128  77.4874ms  78.6021ms
// 256  74.7013ms  76.9963ms
// 512  64.4254ms  78.3839ms
// 1024  62.1404ms  77.6786ms
// 2048  77.4568ms  61.3854ms
// 4096  62.475ms  62.992ms
// 8192  31.2396ms  44.9631ms
// 16384  5.1862ms  5.5706ms
// 32768  10.0455ms  11.5771ms
// 65536  22.6883ms  21.2998ms
// 131072  40.7854ms  41.7471ms
// 262144  83.3409ms  82.7738ms
// 524288  164.3664ms  171.4147ms

func Rwmutex(){
producer := func(wg *sync.WaitGroup, l sync.Locker){
	defer wg.Done()
	for i:=5;i>0;i--{
		l.Lock()
		l.Unlock()
		time.Sleep(1) // slow sleeps for one second
	}
}

observer := func(wg *sync.WaitGroup, l sync.Locker){
	defer wg.Done()
	l.Lock()
	defer l.Unlock()
}

test := func(count int, mutex, rwmutex sync.Locker) time.Duration{
	var wg sync.WaitGroup
	wg.Add(count+1)
	beginTestTime := time.Now()
	go producer(&wg, mutex)
	for i:=count;i>0;i--{
		go observer(&wg, rwmutex)
	}

	wg.Wait()
	return time.Since(beginTestTime)
}

tw:= tabwriter.NewWriter(os.Stdout, 0, 1, 2,' ', 0)
defer tw.Flush()

var m sync.RWMutex

fmt.Fprintf(tw, "Readers\tRWMutex\tMutex\n")
for i:=0;i<20;i++{
	count := int(math.Pow(2, float64(i)))
	fmt.Fprintf(tw,
	"%d\t%v\t%v\n",	
	count,
	test(count, &m, m.RLocker()),
	test(count, &m, &m),
)
}
}
