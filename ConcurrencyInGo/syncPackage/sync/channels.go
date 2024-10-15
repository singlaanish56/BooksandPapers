package sync

import (
	"fmt"
	"sync"
	"time"
)

//simple channels, ranging over the values

func SimpleChannels(){
initstream := make(chan int)
go func(){
	defer close(initstream)
	for i:=1;i<=5;i++{
		initstream <- i
	}
}()

for integer := range initstream{
	fmt.Printf("%v ", integer)
}
}


// closing a channel also works as a broadcast same as the cond variable

func ChannelBroadcast(){
	begin := make(chan interface{})
	var wg sync.WaitGroup

	for i:=0;i<5;i++{
		wg.Add(1)
		go func(i int){
			defer wg.Done()
			<- begin
			fmt.Printf("%v has begin\n", i)
		}(i)
	}


	fmt.Println("unblocking the goroutines")
	close(begin)
	wg.Wait()
}

//select statement are blocking until they find a match

func SimpleSelect(){
	start := time.Now()
	c := make(chan interface{})

	go func(){
		time.Sleep(5*time.Second)
		close(c)
	}()

	fmt.Println("Blocking on the read")
	select {
	case <-c:
		fmt.Printf("Unblocked %v laster .\n", time.Since(start))
	}
}


//runtime decides it between the multiple channels by dividing and giving it time equally

func MultipleSelect(){
	c1 := make(chan interface{}); close(c1)
	c2 := make(chan interface{}); close(c2)

	var c1count, c2count int
	for i:=1000;i>0;i--{
		select {
		case <- c1:
			c1count++
		case <- c2:
			c2count++
		}
	}

	fmt.Printf("c1count: %d\nc2count%d\n", c1count, c2count)
}