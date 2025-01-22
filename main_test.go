package main

import (
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

	// Send some test data
	testData := "PING\r\n"
	_, err = conn.Write([]byte(testData))
	if err != nil {
		t.Fatalf("Failed to send data to server: %v", err)
	}

	// Give the server a moment to process the data
	time.Sleep(100 * time.Millisecond)
}
