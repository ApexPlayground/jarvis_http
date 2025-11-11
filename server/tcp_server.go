package server

import (
	"fmt"
	"net"
	"time"
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

	conn.SetReadDeadline(time.Now().Add(20 * time.Second))

	buffer := make([]byte, 1024) // make slice of 1kb
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Print("Error reading byte:", err)
			break
		}
		fmt.Printf("Read byte: %s", buffer[:n])
		conn.Write([]byte("Message received\n"))

		// Reset the timeout after receiving data
		conn.SetReadDeadline(time.Now().Add(20 * time.Second))

	}

}
