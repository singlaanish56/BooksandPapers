package sync

import (
	"testing"
	"net"
	"io"

)

func init(){
	daemonStarted := StartNetworkDaemon()
	daemonStarted.Wait()
}

//the network operation without the pool

// BenchmarkNetworkRequest-12             8        1323682588 ns/op 1.32 secs
// PASS
// ok      github.com/singlaanish56/Books/ConcurrencyInGo/syncPackage/sync 12.243s


//the network operations with the pool
// BenchmarkNetworkRequest-12            26         506091838 ns/op 0.506 secs
// PASS
// ok      github.com/singlaanish56/Books/ConcurrencyInGo/syncPackage/sync 28.389s

func BenchmarkNetworkRequest(b *testing.B){
	for i:=0;i<b.N;i++{
		conn, err := net.Dial("tcp","localhost:8080")
		if err!=nil{
			b.Fatalf("cannot dial host: %v", err)
		}
		if _, err:= io.ReadAll(conn); err!=nil{
			b.Fatalf("cannot read: %v",err)
		}
		conn.Close()
	}
}