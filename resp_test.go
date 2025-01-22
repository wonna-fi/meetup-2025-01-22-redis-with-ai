package main

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func TestParseRESP(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *RESPValue
		wantErr  bool
	}{
		{
			name:     "simple string",
			input:    "+OK\r\n",
			expected: &RESPValue{Type: SimpleString, Str: "OK"},
		},
		{
			name:     "error",
			input:    "-Error message\r\n",
			expected: &RESPValue{Type: Error, Str: "Error message"},
		},
		{
			name:     "integer",
			input:    ":1000\r\n",
			expected: &RESPValue{Type: Integer, Int: 1000},
		},
		{
			name:     "bulk string",
			input:    "$5\r\nhello\r\n",
			expected: &RESPValue{Type: BulkString, Str: "hello"},
		},
		{
			name:     "null bulk string",
			input:    "$-1\r\n",
			expected: &RESPValue{Type: BulkString, IsNull: true},
		},
		{
			name:  "array",
			input: "*2\r\n$4\r\nPING\r\n$5\r\nhello\r\n",
			expected: &RESPValue{
				Type: Array,
				Array: []RESPValue{
					{Type: BulkString, Str: "PING"},
					{Type: BulkString, Str: "hello"},
				},
			},
		},
		{
			name:     "null array",
			input:    "*-1\r\n",
			expected: &RESPValue{Type: Array, IsNull: true},
		},
		{
			name:    "invalid prefix",
			input:   "invalid\r\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(bytes.NewBufferString(tt.input))
			got, err := ParseRESP(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRESP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseRESP() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRESPValueSerialize(t *testing.T) {
	tests := []struct {
		name     string
		value    RESPValue
		expected string
	}{
		{
			name:     "simple string",
			value:    RESPValue{Type: SimpleString, Str: "OK"},
			expected: "+OK\r\n",
		},
		{
			name:     "error",
			value:    RESPValue{Type: Error, Str: "Error message"},
			expected: "-Error message\r\n",
		},
		{
			name:     "integer",
			value:    RESPValue{Type: Integer, Int: 1000},
			expected: ":1000\r\n",
		},
		{
			name:     "bulk string",
			value:    RESPValue{Type: BulkString, Str: "hello"},
			expected: "$5\r\nhello\r\n",
		},
		{
			name:     "null bulk string",
			value:    RESPValue{Type: BulkString, IsNull: true},
			expected: "$-1\r\n",
		},
		{
			name: "array",
			value: RESPValue{
				Type: Array,
				Array: []RESPValue{
					{Type: BulkString, Str: "PING"},
					{Type: BulkString, Str: "hello"},
				},
			},
			expected: "*2\r\n$4\r\nPING\r\n$5\r\nhello\r\n",
		},
		{
			name:     "null array",
			value:    RESPValue{Type: Array, IsNull: true},
			expected: "*-1\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(tt.value.Serialize())
			if got != tt.expected {
				t.Errorf("RESPValue.Serialize() = %v, want %v", got, tt.expected)
			}
		})
	}
}
