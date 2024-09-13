package main

import (
	"fmt"
	"io"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Connection closed to %s\n", conn.LocalAddr().String())
				return
			}
			fmt.Printf("Error reading: %s\n", err.Error())
		}

		resp := []byte{0x00, 0x01, 0x02, 0x03, 0x00, 0x00, 0x00, 0x07}
		_, err = conn.Write(resp)
		if err != nil {
			fmt.Printf("Error writing: %s\n", err.Error())
		}
	}
}
