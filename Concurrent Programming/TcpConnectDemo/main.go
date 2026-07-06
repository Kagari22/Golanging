package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

const (
	tcp     = "tcp"
	address = "127.0.0.1:5678"
)

func TcpServer() {
	listener, err := net.Listen(tcp, address)
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	log.Printf("%s server is listening on %s\n", tcp, listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("Accept from %s\n", conn.RemoteAddr())
		conn.Close()
	}
}

func TcpClient() {
	num := 10
	wg := sync.WaitGroup{}
	wg.Add(num)

	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()

			conn, err := net.Dial(tcp, address)
			if err != nil {
				log.Println(err)
				return
			}
			defer conn.Close()

			log.Printf("Connection is established, client address is %s\n", conn.LocalAddr())
		}()
	}

	wg.Wait()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage:")
		fmt.Println("  go run . server")
		fmt.Println("  go run . client")
		return
	}

	switch os.Args[1] {
	case "server":
		TcpServer()
	case "client":
		TcpClient()
	default:
		fmt.Println("unknown command:", os.Args[1])
		fmt.Println("use: server or client")
	}
}
