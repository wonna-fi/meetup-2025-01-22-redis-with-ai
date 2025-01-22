# Redis-Lite

Redis-Lite is a lightweight Redis clone implemented in Go. It's a single-node, in-memory database that supports concurrency and implements a subset of Redis commands.

## Prerequisites

- Go 1.23 or higher
- Redis CLI (optional, for testing)

## Building

To build the server:

```bash
go build
```

This will create a `redis-lite` executable in the current directory.

## Running the Server

To start the server:

```bash
./redis-lite
```

The server will start listening on port 6379 (default Redis port).

## Running Tests

To run all tests:

```bash
go test -v
```

To run tests with race condition detection:

```bash
go test -v -race
```

To run specific test suites:

```bash
go test -v -run TestStorage    # Run storage tests only
go test -v -run TestServer    # Run server tests only
```

## Supported Commands

The server currently supports the following Redis commands:

### PING
- Usage: `PING [message]`
- Response: Returns PONG if no argument is provided, otherwise returns the message
- Example:
  ```
  > PING
  PONG
  > PING hello
  hello
  ```

### ECHO
- Usage: `ECHO message`
- Response: Returns the message
- Example:
  ```
  > ECHO hello
  hello
  ```

### SET
- Usage: `SET key value`
- Response: Returns OK if successful
- Example:
  ```
  > SET mykey myvalue
  OK
  ```

### GET
- Usage: `GET key`
- Response: Returns the value of key, or nil if the key doesn't exist
- Example:
  ```
  > GET mykey
  myvalue
  > GET nonexistent
  (nil)
  ```

### DEL
- Usage: `DEL key [key ...]`
- Response: Returns the number of keys that were removed
- Example:
  ```
  > SET key1 value1
  OK
  > SET key2 value2
  OK
  > DEL key1 key2 nonexistent
  2
  ```

## Implementation Details

- Thread-safe in-memory storage using Go's `sync.RWMutex`
- RESP (Redis Serialization Protocol) implementation for client-server communication
- Concurrent request handling using goroutines
- Each client connection is handled in a separate goroutine

## Project Structure

- `main.go` - Server implementation and command handling
- `resp.go` - RESP protocol implementation
- `storage.go` - Thread-safe key-value storage implementation
- `*_test.go` - Test files for each component

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 