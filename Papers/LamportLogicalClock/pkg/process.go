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

func (p *Process) incrementClock() {
	p.Clock += rand.Intn(3)
}

func (p *Process) SendMessage(to int, content string) {
	p.incrementClock()
	m := Message{
		From:      p.ID,
		To:        to,
		Timestamp: p.Clock,
		Content:   content,
	}

	p.RequestQueue <- m
	fmt.Printf("[Send] P%d -> P%d | Time: %d | Content: %s\n", p.ID, to, p.Clock, content)

}

func (p *Process) RecieveMessage(msg Message) {
	p.Clock = max(p.Clock, msg.Timestamp) + 1
	fmt.Printf("[Recieve] P%d | Time: %d | Content: %s\n", p.ID, p.Clock, msg.Content)
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
