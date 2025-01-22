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
			name:         "SET without arguments",
			command:      "*1\r\n$3\r\nSET\r\n",
			expected:     "ERR wrong number of arguments for 'SET' command",
			expectedType: Error,
		},
		{
			name:         "SET with key only",
			command:      "*2\r\n$3\r\nSET\r\n$3\r\nkey\r\n",
			expected:     "ERR wrong number of arguments for 'SET' command",
			expectedType: Error,
		},
		{
			name:         "SET key-value",
			command:      "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n",
			expected:     "OK",
			expectedType: SimpleString,
		},
		{
			name:         "GET non-existent key",
			command:      "*2\r\n$3\r\nGET\r\n$8\r\nnotfound\r\n",
			expected:     "",
			expectedType: BulkString,
			// Note: IsNull will be true for this case
		},
		{
			name:         "GET existing key",
			command:      "*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n",
			expected:     "value",
			expectedType: BulkString,
		},
		{
			name:         "DEL existing key",
			command:      "*2\r\n$3\r\nDEL\r\n$3\r\nkey\r\n",
			expected:     "1",
			expectedType: Integer,
		},
		{
			name:         "DEL non-existent key",
			command:      "*2\r\n$3\r\nDEL\r\n$3\r\nkey\r\n",
			expected:     "0",
			expectedType: Integer,
		},
		{
			name:         "GET after DEL",
			command:      "*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n",
			expected:     "",
			expectedType: BulkString,
			// Note: IsNull will be true for this case
		},
		{
			name:         "Unknown command",
			command:      "*1\r\n$7\r\nUNKNOWN\r\n",
			expected:     "ERR unknown command 'UNKNOWN'",
			expectedType: Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new connection for each test
			conn, err := net.Dial(PROTOCOL, fmt.Sprintf(":%d", DEFAULT_PORT))
			if err != nil {
				t.Fatalf("Failed to connect to server: %v", err)
			}
			defer conn.Close()

			reader := bufio.NewReader(conn)

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

			// For null bulk strings, we only check the type and IsNull flag
			if tt.expectedType == BulkString && tt.expected == "" {
				if !resp.IsNull {
					t.Error("Expected null bulk string")
				}
			} else if resp.Type == Integer {
				// For integer responses, convert the expected string to int64
				expected := int64(0)
				if tt.expected == "1" {
					expected = 1
				}
				if resp.Int != expected {
					t.Errorf("Expected %d, got %d", expected, resp.Int)
				}
			} else if resp.Str != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, resp.Str)
			}

			// Sleep a tiny bit to avoid overwhelming the server
			time.Sleep(10 * time.Millisecond)
		})
	}
}
