package main

import (
	"bufio"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestServerStartup(t *testing.T) {
	// Start server in a goroutine
	go main()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Try to connect to the server
	conn, err := net.Dial(PROTOCOL, fmt.Sprintf(":%d", DEFAULT_PORT))
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	tests := []struct {
		name         string
		command      string
		expected     string
		expectedType RESPType
	}{
		{
			name:         "PING without argument",
			command:      "*1\r\n$4\r\nPING\r\n",
			expected:     "PONG",
			expectedType: SimpleString,
		},
		{
			name:         "PING with argument",
			command:      "*2\r\n$4\r\nPING\r\n$5\r\nhello\r\n",
			expected:     "hello",
			expectedType: SimpleString,
		},
		{
			name:         "ECHO without argument",
			command:      "*1\r\n$4\r\nECHO\r\n",
			expected:     "ERR wrong number of arguments for 'ECHO' command",
			expectedType: Error,
		},
		{
			name:         "ECHO with argument",
			command:      "*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n",
			expected:     "hello",
			expectedType: BulkString,
		},
		{
			name:         "ECHO with too many arguments",
			command:      "*3\r\n$4\r\nECHO\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			expected:     "ERR wrong number of arguments for 'ECHO' command",
			expectedType: Error,
		},
		{
			name:         "Unknown command",
			command:      "*1\r\n$7\r\nUNKNOWN\r\n",
			expected:     "ERR unknown command 'UNKNOWN'",
			expectedType: Error,
		},
	}

	reader := bufio.NewReader(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Send command
			_, err = conn.Write([]byte(tt.command))
			if err != nil {
				t.Fatalf("Failed to send command: %v", err)
			}

			// Read response
			resp, err := ParseRESP(reader)
			if err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			// Verify response type and content
			if resp.Type != tt.expectedType {
				t.Errorf("Expected response type %v, got %v", tt.expectedType, resp.Type)
			}

			if resp.Str != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, resp.Str)
			}
		})
	}
}
