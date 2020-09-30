package examples

import (
	"fmt"
	"testing"
)

func TestBound(t *testing.T) {
	// data はメイン goroutine からもアクセスできるが、規約により loopData からのみアクセスすることにしている
	data := make([]int, 4)

	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- data[i]
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
}

func TestLexical(t *testing.T) {

}
