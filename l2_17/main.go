package main

// telnet util

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type ConnectionParameters struct {
	Host    string
	Port    string
	Timeout time.Duration
}

func main() {
	//
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//
	//// to catch interrupt signals
	//sigChan := make(chan os.Signal, 1)
	//signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	signal.Ignore(syscall.SIGINT, os.Interrupt)
	log.Println("Ctrl+C игнорируется. Используйте Ctrl+D для завершения программы.")

	//go func() {
	//	for sig := range sigChan {
	//		if sig == syscall.SIGINT {
	//			log.Println("Ctrl+C command is handled")
	//		}
	//	}
	//}()

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

	dialer := net.Dialer{Timeout: connectionParameters.Timeout}
	connection, err := dialer.Dial("tcp", remoteAddres)
	if err != nil {
		log.Println(err)
	}

	defer connection.Close()
	log.Println("connection is set")

	//ticker := time.NewTicker(time.Second)
	//secondsCounter := 0
	//isConnected := false
	//for secondsCounter < timeout && !isConnected {
	//	select {
	//	case <-ticker.C:
	//		secondsCounter++
	//		var err error
	//		log.Printf("try to connect number: %d", secondsCounter)
	//		connection, err = net.Dial("tcp", remoteAddres)
	//		if err != nil {
	//			log.Fatalf("tcp connection error: %v", err)
	//		} else {
	//			isConnected = true
	//			log.Println("Connection is set")
	//		}
	//	}
	//}
	//ticker.Stop()
	//defer func() {
	//	err := connection.Close()
	//	if err != nil {
	//		panic(err)
	//	}
	//}()

	// launch two goroutins for writing and reading
	var wg sync.WaitGroup
	var signalChannel = make(chan struct{}, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(signalChannel)
		//scanner := bufio.NewScanner(os.Stdin)
		reader := bufio.NewReader(os.Stdin)
		for {
			input, err := reader.ReadString('\n')

			if err == io.EOF {
				return
			}

			if err != nil {
				log.Println(err)
				//return
			}

			_, err = connection.Write([]byte(input))
			if err != nil {
				log.Printf("tcp write error: %v", err)
				//return
			}
		}
	}()

	//time.Sleep(5 * time.Second)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			buffer := make([]byte, 255)
			select {
			case <-signalChannel:
				log.Println("reading is closed")
				return
			default:
			}

			err := connection.SetDeadline(time.Now().Add(connectionParameters.Timeout))
			if err != nil {
				log.Printf("tcp set connection deadline error: %v", err)
			}

			n, err := connection.Read(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}

				if err == io.EOF {
					log.Printf("tcp read error: %v", err)
				}
				return
			}

			if n > 0 {
				fmt.Print(string(buffer[:n]))
			}
		}
	}()

	wg.Wait()
}
