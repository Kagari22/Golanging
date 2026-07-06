package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		panic(err)
	}
	_, err = conn.Write([]byte("hello tcp"))
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}

	fmt.Println("server response:", string(buf[:n]))
}
