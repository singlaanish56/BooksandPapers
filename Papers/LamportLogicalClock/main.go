package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/singlaanish56/BooksandPapers/pkg"
)

func main() {
	numProcess := 3
	runFor := 10

	var wg sync.WaitGroup

	processes := make([]*pkg.Process, numProcess)
	messageChannels := make([]chan pkg.Message, numProcess)

	for i := 0; i < numProcess; i++ {
		messageChannels[i]= make(chan pkg.Message, 1000)
		processes[i] = pkg.NewProcess(i, messageChannels[i])
		wg.Add(1)
		go processes[i].Run(&wg)
	}

	for runFor > 0 {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		senderId := rand.Intn(numProcess)
		receiverId := rand.Intn(numProcess)
		for receiverId==senderId{
			receiverId = rand.Intn(numProcess)
		}
		content := fmt.Sprintf("Message From process %d" ,senderId)
		processes[senderId].SendMessage(processes[receiverId], content)
		runFor--
	}

	for i:=0;i<numProcess;i++{
		close(messageChannels[i])
	}

	wg.Wait()

	fmt.Println("The clock for all the processes")
	for i := 0; i < numProcess; i++ {
		fmt.Printf("P%d | Time %d\n", processes[i].ID, processes[i].Clock)
	}
}
