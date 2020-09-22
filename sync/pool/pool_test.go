package pool

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

// daemon が起動完了するまで待つ
func init() {
	daemonStarted := startNetworkDaemon()
	daemonStarted.Wait()
}

func connectToService() interface{} {
	time.Sleep(1 * time.Second)
	return struct{}{}
}

func startNetworkDaemon() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			log.Fatalf("cannot listen: %v\n", err)
		}
		defer server.Close()

		wg.Done() // 起動完了

		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("cannnot accept connection: %v\n", err)
				continue
			}
			connectToService() // 1 秒間待つ
			fmt.Fprintln(conn, "")
			conn.Close()
		}
	}()
	return &wg
}

func BenchmarkNetworkRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			b.Fatalf("cannot dial host: %v", err)
		}
		if _, err := ioutil.ReadAll(conn); err != nil {
			b.Fatalf("cannot read: %v", err)
		}
		conn.Close()
	}
}
