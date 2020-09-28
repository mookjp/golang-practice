package channel

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"testing"
)

func TestReceive(t *testing.T) {
	ch := make(chan string)

	go func() {
		ch <- "hello"
	}()

	println(<-ch)
}

func TestBool(t *testing.T) {
	ch := make(chan string)

	go func() {
		ch <- "hello"
	}()

	salutation, ok := <-ch
	fmt.Printf("(%v): %v\n", ok, salutation)
}

func TestClosed(t *testing.T) {
	ch := make(chan int)

	go func() {
		ch <- 100
	}()

	close(ch)
	salutation, ok := <-ch
	fmt.Printf("(%v): %v\n", ok, salutation)
}

func TestClosed2(t *testing.T) {
	ch := make(chan int)

	close(ch)
	go func() {
		ch <- 100
	}()

	salutation, ok := <-ch
	fmt.Printf("(%v): %v\n", ok, salutation)
}

func TestClosed3(t *testing.T) {
	ch := make(chan int)

	go func() {
		ch <- 100
		close(ch)
	}()

	salutation, ok := <-ch
	fmt.Printf("(%v): %v\n", ok, salutation)
}

func TestClose(t *testing.T) {
	beginCh := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-beginCh // channel から読み込みできるようになる（close される）まで待機する
			fmt.Printf("%v has bigun\n", i)
		}(i)
	}
	fmt.Println("unblocking goroutine...")
	close(beginCh) // close して他の goroutine がブロックをやめるようになる
	wg.Wait()
}

func TestBuf(t *testing.T) {
	var stdBuf bytes.Buffer
	defer stdBuf.WriteTo(os.Stdout) // 終了前にバッファから stdout に出力

	intCh := make(chan int, 4) // キャパシティが 4 の channel を作成
	go func() {
		defer close(intCh)
		defer fmt.Fprintln(&stdBuf, "Producer Done.")
		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdBuf, "Sending: %d\n", i) // stdout バッファに溜める
			intCh <- i
		}
	}()

	for integer := range intCh { // バッファに溜まった int を読み込む
		fmt.Fprintf(&stdBuf, "Received: %v.\n", integer)
	}
}

func TestOwner(t *testing.T) {
	chOwner := func() <-chan int {
		// 結果を 6 つ生成することがわかっているため、
		// goroutine をできる限り早く完了するようにキャパシティが 5 のバッファ付き channel を作成
		resultCh := make(chan int, 5)
		go func() {
			defer close(resultCh)
			for i := 0; i <= 5; i++ {
				resultCh <- i
			}
		}()
		return resultCh
	}

	resultCh := chOwner()
	for result := range resultCh {
		fmt.Printf("received: %d\n", result)
	}
	fmt.Println("done receiving!")
}
