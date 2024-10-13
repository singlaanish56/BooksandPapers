package sync

import (
	"fmt"

	"log"
	"net"
	"sync"

	"time"
)

func SimplePool() {
	myPool := &sync.Pool{
		New : func() interface{}{
			fmt.Println("Creating new instance")
			return struct{}{}
		},
	}

	myPool.Get()
	instance := myPool.Get()
	myPool.Put(instance)
	myPool.Get()
}

// trying to pre add a calculator to the pool, so even though the
// the calculator needs a GB of space, it would take only the available ones
func CalculatorPool(){
var numsCalcCreated int

calcPool := &sync.Pool{
	New: func() interface{}{
		numsCalcCreated+=1;
		mem := make([]byte, 1024)
		return &mem
	},
}

//4KB worth of calculator
calcPool.Put(calcPool.New())
calcPool.Put(calcPool.New())
calcPool.Put(calcPool.New())
calcPool.Put(calcPool.New())

const numWorkers = 1024*1024
var wg sync.WaitGroup
wg.Add(numWorkers)

for i:=numWorkers;i>0;i--{
	go func(){
		defer wg.Done()

		mem := calcPool.Get().(*[]byte)
		defer calcPool.Put(mem)
	}()
}

wg.Wait()
fmt.Printf("NUmber of calculators created %d", numsCalcCreated)
}


func connectToService() interface{}{
	time.Sleep(1*time.Second)
	return struct{}{}
}

func warmSrviceConnCache() * sync.Pool{
	p :=&sync.Pool{
		New :connectToService,
		}
	for i:=0;i<10;i++{
			p.Put(p.New())
	}
	return p
}


func StartNetworkDaemon() *sync.WaitGroup{
	var wg sync.WaitGroup
	wg.Add(1)

	go func(){
		connPool := warmSrviceConnCache()

		server, err := net.Listen("tcp","localhost:8080")
		if err !=nil{
			log.Fatalf("cannot listen: %v", err)
		}
		defer server.Close()

		wg.Done()
		for{
			conn, err := server.Accept()
			if err!=nil{
				log.Printf("cannot accept the connection: %v", err)
				continue
			}
			//connectToService()
			svcConn := connPool.Get()
			fmt.Fprintf(conn,"")
			connPool.Put(svcConn)
			conn.Close()
		}
	}()

	return &wg
}

