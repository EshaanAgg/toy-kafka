package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9092")
	if err != nil {
		fmt.Printf("Error connecting: %s\n", err.Error())
	}

	defer conn.Close()
	msg := []byte{0, 0, 0, 0}
	_, err = conn.Write(msg)
	if err != nil {
		fmt.Printf("Error writing: %s\n", err.Error())
	}

	res := make([]byte, 1024)
	n, err := conn.Read(res)
	if err != nil {
		fmt.Printf("Error reading: %s\n", err.Error())
	}
	fmt.Println(res[:n])
}
