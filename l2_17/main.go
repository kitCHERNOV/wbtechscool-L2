package main

// telnet util

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type ConnectionParameters struct {
	Host    string
	Port    string
	Timeout time.Duration
}

func main() {

	var connectionParameters ConnectionParameters = ConnectionParameters{}

	var timeout int
	flag.IntVar(&timeout, "timeout", 10, "")
	flag.Parse()

	// set timeout
	connectionParameters.Timeout = time.Duration(timeout) * time.Second

	// TODO: check if is correct indexes (Ok)
	// Get Host
	connectionParameters.Host = os.Args[2]
	// Get Port
	connectionParameters.Port = os.Args[3]

	fmt.Println(connectionParameters)
	var (
		remoteAddres = net.JoinHostPort(connectionParameters.Host, connectionParameters.Port)
		connection   net.Conn
	)
	//con := net.DialTCP()

	ticker := time.NewTicker(time.Second)
	secondsCounter := 0
	isConnected := false
	for secondsCounter < timeout && !isConnected {
		select {
		case <-ticker.C:
			secondsCounter++
			var err error
			log.Printf("try to connect number: %d", secondsCounter)
			connection, err = net.Dial("tcp", remoteAddres)
			if err != nil {
				log.Fatalf("tcp connection error: %v", err)
			} else {
				isConnected = true
				log.Println("Connection is set")
			}
		}
	}
	ticker.Stop()
	defer func() {
		err := connection.Close()
		if err != nil {
			panic(err)
		}
	}()

	// launch two goroutins for writing and reading
	var wg sync.WaitGroup
	var signalChannel = make(chan struct{}, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(os.Stdin)
		//scanner := bufio.NewReader(os.Stdin)
		for scanner.Scan() {
			signalChannel <- struct{}{}
			_, err := connection.Write([]byte(scanner.Text()))
			if err != nil {
				log.Fatalf("tcp write error: %v", err)
				return
			}
		}
	}()

	//time.Sleep(5 * time.Second)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			buffer := make([]byte, 255)
			<-signalChannel
			err := connection.SetDeadline(time.Now().Add(connectionParameters.Timeout))
			if err != nil {
				log.Fatalf("tcp set connection deadline error: %v", err)
			}
			_, err = connection.Read(buffer)
			if err != nil {
				log.Fatalf("tcp read error: %v", err)
				break
			}

			fmt.Print(string(buffer))
		}
	}()

	wg.Wait()
}
