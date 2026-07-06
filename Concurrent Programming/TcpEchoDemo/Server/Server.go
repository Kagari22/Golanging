package main

import (
	"fmt"
	"net"
)

func handleConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("client disconnected:", err)
			return
		}

		msg := string(buf[:n])
		fmt.Println("received:", msg)

		_, err = conn.Write([]byte("server received: " + msg))
		if err != nil {
			fmt.Println("write error:", err)
			return
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("TCP server listening on 127.0.0.1:9000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		go handleConn(conn)
	}
}
