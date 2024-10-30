package resp

import (
	"errors"
	"strings"
)

// ParseRESPCommand parses RESP protocol to extract the command.
// RESP data format should be like "*<number of elements>\r\n$<number of bytes>\r\nCOMMAND\r\n..."
// see Redis serialization protocol specification
// https://redis.io/docs/latest/develop/reference/protocol-spec/
func ParseRESPCommand(data []byte) (string, error) {
	if len(data) == 0 || data[0] != '*' {
		return "", errors.New("invalid RESP format")
	}

	// Split the data by RESP delimiter "\r\n"
	parts := strings.Split(string(data), "\r\n")
	if len(parts) < 4 || parts[2] == "" {
		return "", errors.New("incomplete RESP command")
	}

	// Command should be the third item in parts if RESP is formatted correctly
	return parts[2], nil
}
