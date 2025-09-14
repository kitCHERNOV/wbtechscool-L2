package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
)

const (
	address = "pool.ntp.org"
)

func main() {
	time, err := ntp.Time(address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting time: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("time: %v\n", time)
}
