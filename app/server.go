package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
	"strings"
)

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	buffer := make([]byte, 4004)
	_, err = conn.Read(buffer)

	if err != nil {
		fmt.Println("Error reading connection: ", err.Error())
		os.Exit(1)
	}
	request := string(buffer)
	requests := strings.Split(request,"\r\n")
	path := strings.Fields(requests[0])[1]

	if path != "/"{
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}else{
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	}
}
