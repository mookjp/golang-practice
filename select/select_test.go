package _select

import (
	"fmt"
	"testing"
	"time"
)

func TestSelect(t *testing.T) {
	start := time.Now()
	c := make(chan interface{})

	go func() {
		time.Sleep(5 * time.Second)
		close(c)
	}()

	fmt.Println("blocking on read...")
	select {
	case <-c:
		fmt.Printf("unblocked %v latter\n", time.Since(start))
	}
}

func TestMultipleChan(t *testing.T) {
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)

	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}
	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
}

func TestAfterChan(t *testing.T) {
	c := make(<-chan int)
	select {
	case <-c:
	case <-time.After(5 * time.Second):
		fmt.Println("timed out")
	}
}

func TestDoneChan(t *testing.T) {
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	workCounter := 0
loop:
	for {
		select {
		case <-done: // done が close した場合に break する
			break loop
		default: // done が close していない場合は以下に続く
		}
		// simulate work
		workCounter++
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("archived %v cycles of work before signalled to stop\n", workCounter)
}
