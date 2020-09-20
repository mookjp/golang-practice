package sync

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCondOfBook(t *testing.T) {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()        // goroutine 内でもクリティカルセクションに入る
		queue = queue[1:] // キューから取り出す処理
		fmt.Println("Removed from queue")
		c.L.Unlock() // クリティカルセクションを抜ける
		c.Signal()   // Wait している goroutine に何か起きたことを知らせる
	}

	for i := 0; i < 10; i++ {
		c.L.Lock() // クリティカルセクションに入る
		for len(queue) == 2 {
			c.Wait() // Cond.Signal() まで待つ。 main goroutine を一時停止する
		}
		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second) // 1秒後にキューから要素を取り出す goroutine
		c.L.Unlock()                        // クリティカルセクションから抜ける
	}
}

func TestCond2(t *testing.T) {
	c := sync.NewCond(&sync.Mutex{})
	var count int

	sayFunc := func() {
		time.Sleep(1 * time.Second)
		c.L.Lock()
		fmt.Printf("hello %d\n", count)
		c.L.Unlock()
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		c.L.Lock()
		count++
		go sayFunc()
		c.Wait()
		c.L.Unlock()
	}
}

func TestNoCond(t *testing.T) {
	var count int
	said := false

	sayFunc := func() {
		fmt.Printf("hello %d\n", count)
		said = true
	}

	for i := 0; i < 10; i++ {
		count++
		go sayFunc()
		for !said {
			time.Sleep(1 * time.Second)
		}
		said = false
	}
}

func TestBroadCast(t *testing.T) {
	type Button struct {
		Clicked *sync.Cond
	}
	button := Button{
		Clicked: sync.NewCond(&sync.Mutex{}),
	}

	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)
		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait() // broadcast を待つ
			fn()     // clickRegistered.Done が実行される
		}()
		goroutineRunning.Wait() // goroutine の発火まで待つ
	}

	var clickRegistered sync.WaitGroup // subscribe を待つための WaitGroup
	clickRegistered.Add(3)

	subscribe(button.Clicked, func() {
		fmt.Println("Maximizing window.")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Mouse clicked.")
		clickRegistered.Done()
	})

	button.Clicked.Broadcast()
	clickRegistered.Wait() // subscribe 実行を待つ
}
