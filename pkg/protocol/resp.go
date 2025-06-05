package protocol

import (
	"fmt"
	"strings"
)

// ParseRESP parses a RESP message into a slice of strings
func ParseRESP(data []byte) []string {
	lines := strings.Split(string(data), "\r\n")
	var out []string
	for _, l := range lines {
		if len(l) > 0 && l[0] != '*' && l[0] != '$' {
			out = append(out, l)
		}
	}
	return out
}

// FormatSubscribe formats a subscribe message
func FormatSubscribe(topic string) string {
	return fmt.Sprintf("*3\r\n$9\r\nsubscribe\r\n$%d\r\n%s\r\n$1\r\n1\r\n", len(topic), topic)
}

// FormatMessage formats a pub/sub message
func FormatMessage(topic, message string) string {
	return fmt.Sprintf("*3\r\n$7\r\nmessage\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n",
		len(topic), topic, len(message), message)
}

// FormatError formats an error message
func FormatError(message string) string {
	return fmt.Sprintf("-ERR %s\r\n", message)
}

// FormatOK formats an OK message
func FormatOK() string {
	return "+OK\r\n"
}

// FormatInteger formats an integer response
func FormatInteger(n int) string {
	return fmt.Sprintf(":%d\r\n", n)
}

// FormatNoAuth formats a no auth message
func FormatNoAuth() string {
	return "-NOAUTH Authentication required\r\n"
}
