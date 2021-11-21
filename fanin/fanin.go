package fanin

import (
	"log"
	"sync"
)

func Funnel(sources ...<-chan int) <-chan int {
	dst := make(chan int) // The shared output channel
	var wg sync.WaitGroup // Used to close dest when all sources are closed
	wg.Add(len(sources))  // Set size of the WaitGroup

	for _, ch := range sources { // Start goroutine for each source
		go func(c <-chan int) {
			defer wg.Done() // Notify WaitGroup when c closes
			for n := range c {
				log.Printf("Adding value %d to dst", n)
				dst <- n
			}
		}(ch)
	}

	go func() { // Start goroutine to close dst when all sources close
		wg.Wait()
		close(dst)
	}()

	return dst
}
