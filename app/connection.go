package main

import (
	"fmt"
	"io"
	"net"

	"github.com/EshaanAgg/toy-kafka/app/datatypes/request"
	"github.com/EshaanAgg/toy-kafka/app/handlers"
)

func handleConnection(conn net.Conn) {
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

		if len(buf) == 0 {
			fmt.Printf("No data received from %s\n", conn.LocalAddr().String())
			return
		}

		// If failed to handle the data recieved from the client, close the connection.
		// TODO: Handle the errors gracefully and inform the client.
		if !handleData(conn, buf) {
			conn.Close()
			return
		}
	}
}

// Handles the recieved data. Returns true if the data was handled successfully.
func handleData(conn net.Conn, data []byte) bool {
	reqHeader, err := request.NewRequestHeader(data)
	if err != nil {
		fmt.Printf("Error in parsing request header: %s\n", err.Error())
		return false
	}

	api, ok := handlers.RequestKeyMap[reqHeader.APIKey]
	if !ok {
		fmt.Printf("Unknown request type: %d\n", reqHeader.APIKey)
		return false
	}

	req, err := api.NewFn(reqHeader)
	if err != nil {
		fmt.Printf("Error in parsing request body: %s\n", err.Error())
		return false
	}

	res, err := req.Handle()
	if err != nil {
		fmt.Printf("Error in handling request: %s\n", err.Error())
		return false
	}

	_, err = conn.Write(res.Bytes())
	if err != nil {
		fmt.Printf("Error writing response: %s\n", err.Error())
		return false
	}

	return true
}
