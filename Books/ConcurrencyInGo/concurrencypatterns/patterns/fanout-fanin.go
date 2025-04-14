package patterns

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"time"
)

//what is fan in and fan out
//based on the concept of pipeline, where you have the different stages.
// what we can use multiple goroutines for a single stage(maybe the heaviest stage)
// and then collate our out put in for the next stage



func toInt(done <-chan interface{}, valueStream <-chan interface{}) <-chan int {
	intStream := make(chan int)
	go func(){
		defer close(intStream)
		for{
			select{
			case <-done:
				return
			case v:= <-valueStream:
				intStream<-v.(int)
			}
		}
	}()
	
	return intStream
}

func isPrime(n int) bool{
    // Handle edge cases
    if n <= 1 {
        return false
    }
    if n <= 3 {
        return true
    }
    if n%2 == 0 || n%3 == 0 {
        return false
    }

    for i := 2; i < n; i++ {
        if n%i == 0 {
            return false
        }
    }
    return true
}

func primeFinder(done <-chan interface{}, intStream <-chan int) <-chan interface{}{
	primeStream := make(chan interface{})

	go func(){
		defer close(primeStream)
		for{
			select{
			case <-done:
				return
			case num := <-intStream:
				if isPrime(num){
					select{
					case <-done:
						return
					case primeStream<-num:
					}
				}
			}
		}
	}()

	return primeStream
}


// just keeps on repeating call the fn you give it
func repeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{}{
	valueStream := make(chan interface{})
	go func(){
		defer close(valueStream)
		for{
			select{
			case <- done:
				return
			case valueStream <- fn():
			}
		}
	}()

	return valueStream
}

// takes the repeat fn, but only runs the int times, so altho the repeat fn runs infinites
// the pipeline benefit helps it oblu run the n+1 times
func take(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{}{
	takeSteam := make(chan interface{})
	go func(){
		defer close(takeSteam)
		for i:=0;i<num;i++{
			select{
			case <-done:
				return
			case takeSteam <- <-valueStream:
			}
		}
	}()

	return takeSteam
}	

// Primes:
//         354241913
//         280690889
//         170683187
//         208683869
//         368592949
//         108509237
//         337056019
//         147563203
//         453557177
//         170867351
// Search took : 9.8313346s


func NoFanOutWithPipeline(){
	done := make(chan interface{})
	defer close(done)
	start:= time.Now()
	
	
	rand := func() interface{} {
        // Generate a random number between 0 and 50000000
        max := big.NewInt(500000000)
        n, err := rand.Int(rand.Reader, max)
        if err != nil {
            return 0 // Handle error case
        }
        return int(n.Int64())
    }



	randIntStream := toInt(done, repeatFn(done, rand))
	fmt.Println("Primes:")
	for prime := range take(done, primeFinder(done, randIntStream), 10){
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took : %v", time.Since(start))
}

//how to decide on the fanout.
//order dependent and the duration, so the prime finder algo, is order dependent, whether  the num is prime or not
// and it take a lot of time  within the normal pipeline
//so to fanoput, we need to strt multiple versions of that stage

//and we need to fanIn the channels 
func fanIn(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{}{
var wg sync.WaitGroup
multiplexed := make(chan interface{})

multiplex := func(c <-chan interface{}){
	defer wg.Done()
	for i:= range c{
		select{
		case <-done:
			return
		case multiplexed <-i:
		}
	}
}

wg.Add(len(channels))
for _,c := range channels{
	go multiplex(c)
}

go func(){
	wg.Wait()
	close(multiplexed)
}()

return multiplexed
}

// Spinning up 12 prime finders. 
// Primes:
//         19730519
//         122699303
//         164722601
//         176321051
//         9914189
//         192000901
//         63183401
//         229992731
//         55903313
//         318464803
// Search took : 2.4221874s  -------- took 1/4 the time

func FanOutWithPipeline(){
	done := make(chan interface{})
	defer close(done)
	start:= time.Now()
	
	
	rand := func() interface{} {
        // Generate a random number between 0 and 50000000
        max := big.NewInt(500000000)
        n, err := rand.Int(rand.Reader, max)
        if err != nil {
            return 0 // Handle error case
        }
        return int(n.Int64())
    }



	randIntStream := toInt(done, repeatFn(done, rand))

	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders. \n", numFinders)
	finders := make([]<-chan interface{}, numFinders)
	fmt.Println("Primes:")
	for i:=0;i<numFinders;i++{
		finders[i] = primeFinder(done, randIntStream)
	}
	for prime := range take(done, fanIn(done, finders...), 10){
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took : %v", time.Since(start))
}
