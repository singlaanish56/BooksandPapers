package patterns

import "fmt"

// Pipelines works same as the you packing one thingo ver the another
// one functions returns the input for another functions and can be stacked over one another
func SimplePipelines(){
	multiply := func(values []int, mulitplier int) []int{
		result := make([]int, len(values))
		for i, v := range values{
			result[i] = v * mulitplier
		}

		return result
	}

	addition := func(values[] int, add int) []int{
		result := make([]int, len(values))
		for i, v := range values{
			result[i] = v+add
		}

		return result
	}


	ints := []int{1,2,3,4,5}
	for _,v := range addition(multiply(ints, 2), 1){
		fmt.Println(v)
	}
}


//using channels to make pipelines
func ChannelPipelines(){
generator := func(done <-chan interface{}, intergers ...int) <-chan int{
	initStream := make(chan int)
	go func(){
		defer close(initStream)
		for _, i:= range intergers{
			select{
			case <-done:
				return
			case initStream <- i:
			}
		}
	}()

	return initStream
}

multiplier := func(done <-chan interface{},inputChannel <-chan int, multiplier int) <-chan int{
	multipliedStream := make(chan int)
	go func(){
		defer close(multipliedStream)
		for i:= range inputChannel{
			select{
			case <-done:
				return
			case multipliedStream <- i*multiplier:
			}
		}
	}()

	return multipliedStream
}

adder := func(done <-chan interface{}, inputChannel <-chan int,add int) <-chan int{
	addedStream := make(chan int)
	go func(){
		defer close(addedStream)
		for i:= range inputChannel{
			select{
			case <-done:
				return
			case addedStream <- i+add:
			}
		}
	}()

	return addedStream
}

done := make(chan interface{})
defer close(done)
 initStream := generator(done, 1,2,3,4,5)
 pipeline :=  multiplier(done, adder(done, multiplier(done, initStream, 2),3),4)

 for v := range pipeline{
	fmt.Println(v)
 } 
}


// because this depends on the go func and channels for the input , you can pass aroung the single array index
// so the pipeline func only waits for one input