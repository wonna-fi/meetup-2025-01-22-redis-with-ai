package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
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

	reader := bufio.NewReader(conn)

	for {
		// Parse RESP message
		value, err := ParseRESP(reader)
		if err != nil {
			log.Printf("Error reading from connection: %v", err)
			return
		}

		// Handle the command
		if value.Type == Array && len(value.Array) > 0 {
			command := strings.ToUpper(value.Array[0].Str)
			args := value.Array[1:]
			log.Printf("Received command: %s, args: %v", command, args)

			var response *RESPValue

			switch command {
			case "PING":
				if len(args) == 0 {
					response = &RESPValue{
						Type: SimpleString,
						Str:  "PONG",
					}
				} else {
					// Echo back the first argument
					response = &RESPValue{
						Type: SimpleString,
						Str:  args[0].Str,
					}
				}
			default:
				response = &RESPValue{
					Type: Error,
					Str:  fmt.Sprintf("ERR unknown command '%s'", command),
				}
			}

			_, err = conn.Write(response.Serialize())
			if err != nil {
				log.Printf("Error writing response: %v", err)
				return
			}
		}
	}
}
