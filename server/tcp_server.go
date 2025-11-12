package server

import (
	"fmt"
	"net"

	"github.com/apexplayground/jarvis_http/handler"
)

func Start() {
	// accept a new incoming TCP connection on port 8080.
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Print("Error listening", err)
	}
	defer listener.Close()

	fmt.Println("Server running on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Print("Error accepting", err)
			continue //await next connection
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	// close connection at end
	defer conn.Close()

	req, err := handler.ParseHttp(conn)
	if err != nil {
		fmt.Print("error parsing http", err)
		return
	}

	fmt.Println("Method:", req.Method)
	fmt.Println("Path:", req.Path)
	fmt.Println("Version:", req.Version)
	fmt.Println("Headers:", req.Headers)

	body := "hello world"

	response := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\n"+
			"Content-Type: text/plain\r\n"+
			"Content-Length: %d\r\n"+
			"\r\n"+
			"%s",
		len(body), body)

	conn.Write([]byte(response))

}
