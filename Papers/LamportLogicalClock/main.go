package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/singlaanish56/BooksandPapers/pkg"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	numProcess := 5
	runFor := 20

	var wg sync.WaitGroup

	processes := make([]*pkg.Process, numProcess)
	messageChannel := make(chan pkg.Message, 100)

	for i := 0; i < numProcess; i++ {
		processes[i] = pkg.NewProcess(i, messageChannel)
		wg.Add(1)
		go processes[i].Run(&wg)
	}

	for runFor > 0 {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		sender := processes[rand.Intn(numProcess)]
		receiverId := rand.Intn(numProcess)
		content := fmt.Sprintf("Message From process %d" ,sender.ID)
		sender.SendMessage(receiverId, content)
		runFor--
	}

	close(messageChannel)
	wg.Wait()

	fmt.Println("The clock for all the processes")
	for i := 0; i < numProcess; i++ {
		fmt.Printf("P%d | Time %d\n", processes[i].ID, processes[i].Clock)
	}
}
