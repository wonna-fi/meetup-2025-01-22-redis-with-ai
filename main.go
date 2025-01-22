package main

import (
	"fmt"
	"log"
	"net"
)

const (
	DEFAULT_PORT = 6379 // Standard Redis port
	PROTOCOL     = "tcp"
)

func main() {
	listener, err := net.Listen(PROTOCOL, fmt.Sprintf(":%d", DEFAULT_PORT))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Printf("Redis-lite server listening on port %d", DEFAULT_PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Printf("New connection from %s", conn.RemoteAddr())

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Connection closed from %s: %v", conn.RemoteAddr(), err)
			return
		}

		log.Printf("Received %d bytes from %s: %s", n, conn.RemoteAddr(), buffer[:n])
	}
}
