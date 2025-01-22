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

var storage *Storage

func main() {
	storage = NewStorage()

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
			case "ECHO":
				if len(args) != 1 {
					response = &RESPValue{
						Type: Error,
						Str:  "ERR wrong number of arguments for 'ECHO' command",
					}
				} else {
					response = &RESPValue{
						Type: BulkString,
						Str:  args[0].Str,
					}
				}
			case "SET":
				if len(args) != 2 {
					response = &RESPValue{
						Type: Error,
						Str:  "ERR wrong number of arguments for 'SET' command",
					}
				} else {
					storage.Set(args[0].Str, args[1].Str)
					response = &RESPValue{
						Type: SimpleString,
						Str:  "OK",
					}
				}
			case "GET":
				if len(args) != 1 {
					response = &RESPValue{
						Type: Error,
						Str:  "ERR wrong number of arguments for 'GET' command",
					}
				} else {
					if value, exists := storage.Get(args[0].Str); exists {
						log.Printf("GET %s: found value %s", args[0].Str, value)
						response = &RESPValue{
							Type: BulkString,
							Str:  value,
						}
					} else {
						log.Printf("GET %s: key not found", args[0].Str)
						response = &RESPValue{
							Type:   BulkString,
							IsNull: true,
						}
					}
				}
			case "DEL":
				if len(args) < 1 {
					response = &RESPValue{
						Type: Error,
						Str:  "ERR wrong number of arguments for 'DEL' command",
					}
				} else {
					// Extract keys from arguments
					keys := make([]string, len(args))
					for i, arg := range args {
						keys[i] = arg.Str
					}

					// Delete keys and get count of deleted keys
					deleted := storage.Del(keys...)
					response = &RESPValue{
						Type: Integer,
						Int:  deleted,
					}
				}
			default:
				response = &RESPValue{
					Type: Error,
					Str:  fmt.Sprintf("ERR unknown command '%s'", command),
				}
			}

			log.Printf("Sending response: %v", response)
			_, err = conn.Write(response.Serialize())
			if err != nil {
				log.Printf("Error writing response: %v", err)
				return
			}
		}
	}
}
