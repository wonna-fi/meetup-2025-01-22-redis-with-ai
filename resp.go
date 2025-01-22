package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	// RESP type prefixes
	ArrayPrefix  = '*'
	BulkPrefix   = '$'
	StringPrefix = '+'
	ErrorPrefix  = '-'
	IntPrefix    = ':'

	// Protocol delimiters
	CR = '\r'
	LF = '\n'
)

// RESPType represents different RESP data types
type RESPType int

const (
	SimpleString RESPType = iota
	Error
	Integer
	BulkString
	Array
)

// RESPValue represents a RESP protocol value
type RESPValue struct {
	Type   RESPType
	Str    string
	Int    int64
	Array  []RESPValue
	IsNull bool
}

// ParseRESP parses a RESP message from a reader
func ParseRESP(reader *bufio.Reader) (*RESPValue, error) {
	// Read the first byte to determine the type
	prefix, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch prefix {
	case ArrayPrefix:
		return parseArray(reader)
	case BulkPrefix:
		return parseBulkString(reader)
	case StringPrefix:
		return parseSimpleString(reader)
	case ErrorPrefix:
		return parseError(reader)
	case IntPrefix:
		return parseInteger(reader)
	default:
		return nil, fmt.Errorf("unknown RESP type prefix: %c", prefix)
	}
}

// readLine reads until CR LF and returns the line without CR LF
func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString(LF)
	if err != nil {
		return "", err
	}
	if len(line) < 2 || line[len(line)-2] != CR {
		return "", fmt.Errorf("invalid RESP line ending")
	}
	return line[:len(line)-2], nil // Remove CR LF
}

func parseArray(reader *bufio.Reader) (*RESPValue, error) {
	// Read array length
	line, err := readLine(reader)
	if err != nil {
		return nil, err
	}

	length, err := strconv.ParseInt(line, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid array length: %s", line)
	}

	if length == -1 {
		return &RESPValue{Type: Array, IsNull: true}, nil
	}

	// Parse array elements
	array := make([]RESPValue, length)
	for i := int64(0); i < length; i++ {
		value, err := ParseRESP(reader)
		if err != nil {
			return nil, err
		}
		array[i] = *value
	}

	return &RESPValue{Type: Array, Array: array}, nil
}

func parseBulkString(reader *bufio.Reader) (*RESPValue, error) {
	// Read string length
	line, err := readLine(reader)
	if err != nil {
		return nil, err
	}

	length, err := strconv.ParseInt(line, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid bulk string length: %s", line)
	}

	if length == -1 {
		return &RESPValue{Type: BulkString, IsNull: true}, nil
	}

	// Read the string content
	data := make([]byte, length)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return nil, err
	}

	// Read and verify CR LF
	cr, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	lf, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	if cr != CR || lf != LF {
		return nil, fmt.Errorf("invalid bulk string ending")
	}

	return &RESPValue{Type: BulkString, Str: string(data)}, nil
}

func parseSimpleString(reader *bufio.Reader) (*RESPValue, error) {
	line, err := readLine(reader)
	if err != nil {
		return nil, err
	}
	return &RESPValue{Type: SimpleString, Str: line}, nil
}

func parseError(reader *bufio.Reader) (*RESPValue, error) {
	line, err := readLine(reader)
	if err != nil {
		return nil, err
	}
	return &RESPValue{Type: Error, Str: line}, nil
}

func parseInteger(reader *bufio.Reader) (*RESPValue, error) {
	line, err := readLine(reader)
	if err != nil {
		return nil, err
	}

	n, err := strconv.ParseInt(line, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid integer: %s", line)
	}
	return &RESPValue{Type: Integer, Int: n}, nil
}

// Serialize converts a RESPValue to its wire format
func (v *RESPValue) Serialize() []byte {
	switch v.Type {
	case SimpleString:
		return []byte(fmt.Sprintf("+%s\r\n", v.Str))
	case Error:
		return []byte(fmt.Sprintf("-%s\r\n", v.Str))
	case Integer:
		return []byte(fmt.Sprintf(":%d\r\n", v.Int))
	case BulkString:
		if v.IsNull {
			return []byte("$-1\r\n")
		}
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v.Str), v.Str))
	case Array:
		if v.IsNull {
			return []byte("*-1\r\n")
		}
		result := []byte(fmt.Sprintf("*%d\r\n", len(v.Array)))
		for _, item := range v.Array {
			result = append(result, item.Serialize()...)
		}
		return result
	default:
		return []byte(fmt.Sprintf("-ERR unknown value type\r\n"))
	}
}
