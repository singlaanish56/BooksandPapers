package pkg

import (
	"fmt"
	"math/rand"
	"sync"
)

type Process struct {
	ID           int
	Clock        int
	RequestQueue chan Message
	ProcessMutex  sync.Mutex 
}

func NewProcess(id int, msgChan chan Message) *Process {
	return &Process{
		ID:           id,
		Clock:        0,
		RequestQueue: msgChan,
	}
}

func (p *Process) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		msg, ok := <-p.RequestQueue
		if !ok {
			return
		}
		p.RecieveMessage(msg)
	}
}

func (p *Process) SendMessage(to *Process, content string) {
	p.ProcessMutex.Lock()
	defer p.ProcessMutex.Unlock()

	p.Clock += rand.Intn(3)

	m := Message{
		From:      p.ID,
		To:        to.ID,
		Timestamp: p.Clock,
		Content:   content,
	}

	to.RequestQueue <- m
	fmt.Printf("[Send] P%d -> P%d | Time: %d | Content: %s\n", p.ID, to.ID, p.Clock, content)

}

func (p *Process) RecieveMessage(msg Message) {
	p.ProcessMutex.Lock()
	defer p.ProcessMutex.Unlock()
	p.Clock = max(p.Clock, msg.Timestamp) + 1
	fmt.Printf("[Recieve] P%d from P%d | Time: %d | Content: %s\n", p.ID, msg.From, p.Clock, msg.Content)
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
