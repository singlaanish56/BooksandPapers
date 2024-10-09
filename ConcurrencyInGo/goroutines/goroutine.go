package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	//this one prints welcome, so salutation operates on the same memory as the main thread

	// salutation := "hello"
	// wg.Add(1)

	
	// go func(){
	// 	defer wg.Done()
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
}