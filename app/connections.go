package main

import (
	"fmt"
	"io"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/protocol"
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

		req := protocol.NewRequest(buf)

		res := protocol.NewResponse()
		res.WriteInt32(0)
		res.WriteInt32(req.Header.CorrelationID)
		res.Send(&conn)
	}
}
