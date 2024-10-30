package resp_test

import (
	"testing"

	"github.com/zombocoder/go-redis-proxy/pkg/resp"
)

func TestParseRESPCommand(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expected    string
		expectError bool
	}{
		{
			name:        "Valid GET command",
			input:       []byte("*2\r\n$3\r\nGET\r\n$8\r\nmykey\r\n"),
			expected:    "GET",
			expectError: false,
		},
		{
			name:        "Valid SET command",
			input:       []byte("*3\r\n$3\r\nSET\r\n$8\r\nmykey\r\n$5\r\nvalue\r\n"),
			expected:    "SET",
			expectError: false,
		},
		{
			name:        "Empty data",
			input:       []byte(""),
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid format - missing *",
			input:       []byte("2\r\n$3\r\nGET\r\n$8\r\nmykey\r\n"),
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid command - missing command",
			input:       []byte("*2\r\n$3\r\n\r\n$8\r\nmykey\r\n"),
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, err := resp.ParseRESPCommand(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got none")
				} else {
					t.Logf("received expected error: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else if command != tt.expected {
					t.Errorf("expected command %q, got %q", tt.expected, command)
				}
			}
		})
	}
}
