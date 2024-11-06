package patterns

import (
	"bytes"
	"fmt"
	"sync"
)

// confinement ensures that the information is only ever
//available from one concurrent process

//adhoc confinement
// as the name suggest this confinment is ensured inherently by how the code
//is strucutred,
// this type is  heavily dependent on discipline to maintain the confinement.

// in this case the access to the data is available both to the loopData function
//and also to the handData loop, this dependent on a a developers mistake to ruin this type of confinement
func AdHocConfinement(){
	data := make([]int, 4)

	loopData := func(handleData chan<- int){
		defer close(handleData)
		for i:= range data{
			handleData<-data[i]
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData{
		fmt.Println(num)
	}
}

//lexical scope
// in this scenario the access to the data is limited to a specific block in the code
// ans hence its easier to enforce, because consumer -> consumes, producer -> produces

func LexicalConfinement(){

	chanOwner := func() <-chan int{
		results := make(chan int, 5)
		go func(){
			defer close(results)
			for  i:=0;i<5;i++{
				results <- i
			}
		}()
		return results
	}

	consumer := func(results <-chan int){
		for result := range results{
			fmt.Printf("Recieved: %d\n", result)
		}

		fmt.Println("Done recieving!")
	}

	results := chanOwner()
	consumer(results)

}

//altho channels are inherently thread safe
//another way to use the lexical confinement for DS which is not a concurrent safe

// in this case altho the printData function operats on different part of the data slice
// but the data slice in itself is not thread safe, in this case the slice should rather be copied
func LexicalConfinement2(){

	printData := func(wg * sync.WaitGroup, data []byte){
		defer wg.Done()

		var buff bytes.Buffer
		for _, b:= range data{	
			fmt.Fprintf(&buff, "%c", b)
		}

		fmt.Println(buff.String())
	}

	var wg sync.WaitGroup

	wg.Add(2)
	data := []byte("golang")
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])

	wg.Wait()

}