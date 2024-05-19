package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
	"strings"
	"encoding/json"
)
type Response struct{
	Version string
	Status_code int
	Status_string string
	Content_type string `json:"Content-Type,omitempty"`
	Content_length int `json:"Content-Length,omitempty"`
	Body string `json:"Body,omitempty"`
}
type Request struct {
    Method string
	Path string
	Version string
	Host string
	Accept string
	User_agent string `json:"User-Agent"`
	Body string
}

func (r *Request) SetRequest(request_lines []string) {
    start_line := strings.Split(request_lines[0]," ")
	request_lines = request_lines[1:]
	dict := make(map[string]interface{})
	dict["method"] = start_line[0]
	dict["path"] = start_line[1]
	dict["version"] = start_line[2]

    for _, value := range request_lines {
		request_line_split := strings.Split(value,": ")

		if len(request_line_split) > 1{
			dict[request_line_split[0]] = request_line_split[1]
		} 
    }
	jsonbody, err := json.Marshal(dict)
	if err != nil {
        fmt.Println(err)
        return
    }
	
    if err := json.Unmarshal(jsonbody, &r); err != nil {
        // do error check
        fmt.Println(err)
        return
    }
}

func (r *Response) SetRespond(Status_code int, Status_string string,options Response) {
	r.Version = "HTTP/1.1"
	r.Status_code = Status_code
	r.Status_string = Status_string
	r.Content_type = options.Content_type 
	r.Content_length = len([]byte(options.Body))
	r.Body = options.Body
}

func (r *Response) SendRespond(conn net.Conn) {
	responseData, err := json.Marshal(r)
    if err!= nil {
		panic(err)
    }
	dict := make(map[string]interface{})

    err = json.Unmarshal([]byte(responseData), &dict)

    if err != nil {
        panic(err)
    }
	response_string := fmt.Sprintf("%s %d %s\r\n", r.Version, r.Status_code, r.Status_string)
	delete(dict, "Status_code")
	delete(dict, "Status_string")
	delete(dict, "Version")
	for key, val := range dict {
		if key == "Body"{
			line := fmt.Sprintf("\r\n%v", val)
			response_string = response_string + line

		}else{
			line := fmt.Sprintf("%s: %v\r\n",key, val)
			response_string = response_string + line
		}
	}
	
	fmt.Printf(fmt.Sprintf(response_string))
	conn.Write([]byte(response_string))
}



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
	request_lines := strings.Split(request,"\r\n")
	r := new(Request)
	r.SetRequest(request_lines)

	response := new(Response)
	path := r.Path
	switch {
		case path == "/":
			response.SetRespond(200,"OK",Response{})

		case strings.Contains(path, "/echo"):
			path_string := strings.TrimPrefix(path, "/echo/")
			response.SetRespond(200,"OK",Response{Body: path_string})
		
		case path == "/user-agent":
			response.SetRespond(200,"OK",Response{Content_type:"text/plain",Body: r.User_agent})

		default:
			response.SetRespond(404,"Not Found",Response{})
	}
	response.SendRespond(conn)
}



