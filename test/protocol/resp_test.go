package protocol_test

import (
	"testing"

	"redix/pkg/protocol"
)

func TestParseRESP(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []string
	}{
		{
			name:     "simple command",
			input:    []byte("*1\r\n$4\r\nPING\r\n"),
			expected: []string{"PING"},
		},
		{
			name:     "command with arguments",
			input:    []byte("*2\r\n$4\r\nAUTH\r\n$8\r\ntoken123\r\n"),
			expected: []string{"AUTH", "token123"},
		},
		{
			name:     "empty input",
			input:    []byte(""),
			expected: []string{},
		},
		{
			name:     "complex RESP array",
			input:    []byte("*3\r\n$9\r\nsubscribe\r\n$4\r\ntest\r\n$1\r\n1\r\n"),
			expected: []string{"subscribe", "test", "1"},
		},
		{
			name:     "multiple commands",
			input:    []byte("*2\r\n$4\r\nAUTH\r\n$8\r\ntoken123\r\n*1\r\n$4\r\nPING\r\n"),
			expected: []string{"AUTH", "token123", "PING"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := protocol.ParseRESP(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("ParseRESP() length = %v, want %v", len(result), len(tt.expected))
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("ParseRESP()[%d] = %v, want %v", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestFormatSubscribe(t *testing.T) {
	tests := []struct {
		name     string
		topic    string
		expected string
	}{
		{
			name:     "simple topic",
			topic:    "test",
			expected: "*3\r\n$9\r\nsubscribe\r\n$4\r\ntest\r\n$1\r\n1\r\n",
		},
		{
			name:     "empty topic",
			topic:    "",
			expected: "*3\r\n$9\r\nsubscribe\r\n$0\r\n\r\n$1\r\n1\r\n",
		},
		{
			name:     "long topic",
			topic:    "very/long/topic/name",
			expected: "*3\r\n$9\r\nsubscribe\r\n$20\r\nvery/long/topic/name\r\n$1\r\n1\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := protocol.FormatSubscribe(tt.topic)
			if result != tt.expected {
				t.Errorf("FormatSubscribe() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFormatMessage(t *testing.T) {
	tests := []struct {
		name     string
		topic    string
		message  string
		expected string
	}{
		{
			name:     "simple message",
			topic:    "test",
			message:  "hello",
			expected: "*3\r\n$7\r\nmessage\r\n$4\r\ntest\r\n$5\r\nhello\r\n",
		},
		{
			name:     "empty message",
			topic:    "test",
			message:  "",
			expected: "*3\r\n$7\r\nmessage\r\n$4\r\ntest\r\n$0\r\n\r\n",
		},
		{
			name:     "long message",
			topic:    "test",
			message:  "this is a very long message that needs to be tested",
			expected: "*3\r\n$7\r\nmessage\r\n$4\r\ntest\r\n$51\r\nthis is a very long message that needs to be tested\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := protocol.FormatMessage(tt.topic, tt.message)
			if result != tt.expected {
				t.Errorf("FormatMessage() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFormatError(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "simple error",
			message:  "invalid token",
			expected: "-ERR invalid token\r\n",
		},
		{
			name:     "empty error",
			message:  "",
			expected: "-ERR \r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := protocol.FormatError(tt.message)
			if result != tt.expected {
				t.Errorf("FormatError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFormatOK(t *testing.T) {
	expected := "+OK\r\n"
	result := protocol.FormatOK()
	if result != expected {
		t.Errorf("FormatOK() = %v, want %v", result, expected)
	}
}

func TestFormatInteger(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected string
	}{
		{
			name:     "positive number",
			n:        42,
			expected: ":42\r\n",
		},
		{
			name:     "zero",
			n:        0,
			expected: ":0\r\n",
		},
		{
			name:     "negative number",
			n:        -1,
			expected: ":-1\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := protocol.FormatInteger(tt.n)
			if result != tt.expected {
				t.Errorf("FormatInteger() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFormatNoAuth(t *testing.T) {
	expected := "-NOAUTH Authentication required\r\n"
	result := protocol.FormatNoAuth()
	if result != expected {
		t.Errorf("FormatNoAuth() = %v, want %v", result, expected)
	}
}
