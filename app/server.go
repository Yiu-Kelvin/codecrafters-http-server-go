package main

import (
	"fmt"
	// Uncomment this block to passddd the first stage
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

	if path == "/"{
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	}else if strings.HasPrefix(path, "/echo"){
		path_string := strings.TrimPrefix(path, "/echo/")
		path_string_len := len([]byte(path_string))
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", path_string_len, path_string)))
	}else if strings.HasPrefix(path, "/user-agent"){
		user_agent := strings.Fields(requests[3])[1]
		user_agent_len := len([]byte(user_agent))
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", user_agent_len, user_agent)))

	}else{
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}



