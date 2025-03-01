package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/singlaanish56/BooksandPapers/pkg"
)


func main(){

	numberOfProcess := 3
	runFor:=5
	processes := make([]*pkg.Process ,numberOfProcess)
	messageChannels := make([]chan pkg.Message, numberOfProcess)

	var wg sync.WaitGroup
	for i:=0;i<numberOfProcess;i++{
		messageChannels[i] = make(chan pkg.Message, 1000)
		processes[i] = pkg.NewProcess(i, messageChannels[i])

		wg.Add(1)
		go processes[i].Run(&wg)
	}

	for runFor >0{
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		senderID := rand.Intn(numberOfProcess)
		receiverID := rand.Intn(numberOfProcess)
		for receiverID==senderID{
			receiverID = rand.Intn(numberOfProcess)
		}

		content := fmt.Sprintf("This is a message from process %d\n", senderID)
		processes[senderID].SendMessage(processes[receiverID], content)
		runFor --;
	}

	for i:=0;i<numberOfProcess;i++{
		close(messageChannels[i])
	}

	wg.Wait()

	fmt.Println("the vector Clocks for various process")
	for i:=0;i<numberOfProcess;i++{
		fmt.Printf("P%d | T ", i);
		fmt.Println(processes[i].VectorClock)
	}

}