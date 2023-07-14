package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
)

func TestHandleConnection(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{
			name:           "SET command",
			input:          "SET key value\n",
			expectedOutput: "OK",
		},
		{
			name:           "SET command with quotes",
			input:          "SET \"key\" \"value\"\n",
			expectedOutput: "OK",
		},
		{
			name:           "GET command",
			input:          "GET key\n",
			expectedOutput: "value",
		},
		{
			name:           "DEL command",
			input:          "DEL key\n",
			expectedOutput: "1",
		},
		{
			name:           "Unknown command",
			input:          "UNKNOWN command\n",
			expectedOutput: "(error) ERR unknown command `UNKNOWN`, with args beginning with: `command`,",
		},
		{
			name:           "SET command with unbalanced quotes",
			input:          "SET \"key\" \"value\n",
			expectedOutput: "(error) ERR Protocol error: unbalanced quotes in request",
		},
	}

	// Start the server in a separate goroutine
	go func() {
		main()
	}()

	for _, testCase := range tests {
		// Connect to the server
		conn, err := net.Dial("tcp", "localhost:9736")
		if err != nil {
			t.Fatalf("Failed to connect to server: %v", err)
		}
		defer conn.Close()

		// Send the test message to the server
		_, err = fmt.Fprintf(conn, testCase.input)
		if err != nil {
			t.Fatalf("Failed to send message to server: %v", err)
		}

		// Read the response from the server
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read response from server: %v", err)
		}

		response = strings.Trim(response, "$")
		response = strings.Trim(response, "\n")

		// Compare the received response with the expected response
		if response != testCase.expectedOutput {
			t.Errorf("Expected response: %s, but got: %s", testCase.expectedOutput, response)
		}
	}
}
