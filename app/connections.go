package main

import (
	"fmt"
	"net"
)

func handleConnection(conn *net.Conn) {
	for {
		buf := make([]byte, 1024)
		_, err := (*conn).Read(buf)
		if err != nil {
			fmt.Printf("Error reading: %s\n", err.Error())
		}

		resp := []byte{0x00, 0x01, 0x02, 0x03, 0x00, 0x00, 0x00, 0x07}
		_, err = (*conn).Write(resp)
		if err != nil {
			fmt.Printf("Error writing: %s\n", err.Error())
		}
	}
}
