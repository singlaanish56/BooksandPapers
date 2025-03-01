package pkg

import (
	"fmt"
	"math/rand"
	"sync"
)

type Process struct{
	Id int
	VectorClock map[int]int
	Request chan Message
	VectorMutex sync.Mutex
}


func NewProcess(id int , msgChan chan Message) *Process{
	p := Process{
		Id : id,
		VectorClock: make(map[int]int),
		Request: msgChan,
	}

	p.VectorClock[id]=0
	return &p
}

func (p *Process) Run(wg *sync.WaitGroup){
	defer wg.Done()
	for{
		msg, ok := <- p.Request
		if !ok{
			return
		}
		p.RecieveMessage(msg)
	}
}


func (p * Process) SendMessage(to *Process, content string){
	p.VectorMutex.Lock()
	defer p.VectorMutex.Unlock()
	p.VectorClock[p.Id] += rand.Intn(3)+1

    timestampCopy := make(map[int]int)
    for k, v := range p.VectorClock {
        timestampCopy[k] = v
    }
    
    m := Message{
        SenderId:   p.Id,
        ReceiverID: to.Id,
        Timestamp:  timestampCopy,
        Content:    content,
    }
	
    to.Request <- m
	fmt.Printf("[Send] P%d -> P%d| Time: %v | Content: %s\n", p.Id, to.Id, p.VectorClock, content)
}

func max(a, b int) int{
	if(a > b){
		return a
	}

	return b
}

func (p *Process) RecieveMessage(msg Message){
	p.VectorMutex.Lock()
	defer p.VectorMutex.Unlock()
	//find if the sender's process exist in the timestamp
	v, ok := p.VectorClock[msg.SenderId]
	if (ok && v<=msg.Timestamp[msg.SenderId]) || !ok{
		p.VectorClock[msg.SenderId] = 1+msg.Timestamp[msg.SenderId]
	}

	for k, v1 := range msg.Timestamp{
		v2, ok := p.VectorClock[k];
		if ok{
			p.VectorClock[k] = max(v1, v2)
		}else{
			p.VectorClock[k]=v1
		}
	}

	fmt.Printf("[Recieve] P%d from P%d | Updated Timestamp: %v | Content: %s \n", p.Id, msg.SenderId, p.VectorClock, msg.Content)
}
