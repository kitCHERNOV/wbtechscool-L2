package main

import (
	"fmt"
	"time"
)

func or(channels ...<-chan interface{}) <-chan struct{} {
	// Both done channel
	var bothDoneChannel = make(chan struct{})
	// Synchronised channel ato get any signal from any channel
	var syncChan = make(chan interface{}, 1)

	// Create func to monitor input channels
	monitor := func(ch <-chan interface{}) {
		syncChan <- <-ch
	}

	// Get any signal from any monitor goroutines
	for _, channel := range channels {
		go monitor(channel)
	}

	// Wait for any signal from goroutines pool
	go func() {
		// Close the bothDoneChannel when any of workers send the signal to channel
		<-syncChan
		close(bothDoneChannel)
	}()

	return bothDoneChannel
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}
