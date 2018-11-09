package main

import (
	"fmt"
	"sync"

	"github.com/jonstacks/patterns/pkg/broadcast"
)

func processChan(name string, c <-chan string) {
	for msg := range c {
		fmt.Printf("%s: %s\n", name, msg)
	}
}

func main() {
	cin := make(chan string)
	cout1 := make(chan string, 5)
	cout2 := make(chan string, 5)

	// Start the broadcasting from input to outputs. Close output channels when
	// input channel is closed
	broadcast.Strings(cin, []chan string{cout1, cout2}, true)

	cin <- "Hello"
	cin <- "World"
	close(cin)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer func() { wg.Done() }()
		processChan("cout1", cout1)
	}()

	wg.Add(1)
	go func() {
		defer func() { wg.Done() }()
		processChan("cout2", cout2)
	}()

	wg.Wait()
}
