package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	//var wg sync.WaitGroup

	//this one prints welcome, so salutation operates on the same memory as the main thread

	// salutation := "hello"
	// wg.Add(1)

	
	// go func(){
	// 	defer g.Done()
	// 	salutation = "welcome"
	// }()

	// wg.Wait()
	// fmt.Println(salutation)

	//-----------------------------------------------------------------------------------

	//in this case the last element in the slice is printed, why because in this case
	// the go routine is run after the for loop or the main go routine completes its work
	// so the go runtime being  very chalak stores the last element of the saturation variable in the heap, because it knows a go routine is trying to access that variable

	// for _, salutation := range []string{"hello","good day","morning"}{
	// 	wg.Add(1)
	// 	go func(){
	// 		defer wg.Done()
	// 		fmt.Println(salutation)
	// 	}()
	// }
	//wg.Wait()

	//---------------------------------------------------------------------------------

	//capturing the variable in the closure helps in printing all the values of the slice, cause now it operates within the loop itself

	// for _, salutation := range []string{"hello","good day","morning"}{
	// 	wg.Add(1)
	// 	go func(salutation string){
	// 		defer wg.Done()
	// 		fmt.Println(salutation)
	// 	}(salutation)
	// }
	// wg.Wait()


	//-------------------------------------------------------------------------------------	
	
	//forever gor outines are not garbage collected by the goroutine runtime
	// still takes 9kb of space in my system

	// memConsumed := func() uint64 {
	// 	runtime.GC()
	// 	var s runtime.MemStats
	// 	runtime.ReadMemStats(&s)
	// 	return s.Sys
	// }

	// var c <-chan interface{}
	// var wg  sync.WaitGroup
	// noop := func(){wg.Done(); <-c}

	// const numOfGORoutines = 1e4
	// wg.Add(numOfGORoutines)
	// before := memConsumed()
	// for i:=numOfGORoutines;i>0;i--{
	// 	go noop()
	// }

	// wg.Wait()
	// after := memConsumed()
	// fmt.Printf("%.3fkb",float64(after-before)/numOfGORoutines/1000)
}