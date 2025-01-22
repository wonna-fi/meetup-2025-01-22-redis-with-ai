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

	// Send PING command in RESP format
	pingCmd := "*1\r\n$4\r\nPING\r\n"
	_, err = conn.Write([]byte(pingCmd))
	if err != nil {
		t.Fatalf("Failed to send PING command: %v", err)
	}

	// Read response
	reader := bufio.NewReader(conn)
	resp, err := ParseRESP(reader)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	// Verify response
	if resp.Type != SimpleString || resp.Str != "PONG" {
		t.Errorf("Expected PONG response, got %v", resp)
	}
}
