package main

import (
	"fmt"
	"io"
	"net"

	"github.com/EshaanAgg/toy-kafka/app/datatypes/request"
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

		if len(buf) == 0 {
			fmt.Printf("No data received from %s\n", conn.LocalAddr().String())
			return
		}

		handleData(conn, buf)
	}
}

func handleData(conn net.Conn, data []byte) {
	reqHeader, err := request.NewRequestHeader(data)
	if err != nil {
		fmt.Printf("Error in parsing request header: %s\n", err.Error())
		return
	}

	api, ok := request.RequestKeyMap[reqHeader.APIKey]
	if !ok {
		fmt.Printf("Unknown request type: %d\n", reqHeader.APIKey)
		return
	}

	req, err := api.NewFn(reqHeader)
	if err != nil {
		fmt.Printf("Error in parsing request body: %s\n", err.Error())
		return
	}

	res, err := req.Handle()
	if err != nil {
		fmt.Printf("Error in handling request: %s\n", err.Error())
		return
	}

	_, err = conn.Write(res.Bytes())
	if err != nil {
		fmt.Printf("Error writing response: %s\n", err.Error())
		return
	}
}
