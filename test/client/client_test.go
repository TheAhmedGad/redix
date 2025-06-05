package client_test

import (
	"net"
	"testing"
	"time"

	"redix/pkg/client"
)

// mockConn is a mock implementation of net.Conn for testing
type mockConn struct {
	readData  []byte
	writeData []byte
	closed    bool
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	if len(m.readData) == 0 {
		return 0, nil
	}
	n = copy(b, m.readData)
	m.readData = m.readData[n:]
	return n, nil
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	m.writeData = append(m.writeData, b...)
	return len(b), nil
}

func (m *mockConn) Close() error {
	m.closed = true
	return nil
}

func (m *mockConn) LocalAddr() net.Addr  { return nil }
func (m *mockConn) RemoteAddr() net.Addr { return nil }

// These methods are required by net.Conn
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestNew(t *testing.T) {
	conn := &mockConn{}
	c := client.New(conn)

	if c.Conn != conn {
		t.Errorf("New() Conn = %v, want %v", c.Conn, conn)
	}
	if c.Subs == nil {
		t.Error("New() Subs map is nil")
	}
	if c.Authed {
		t.Error("New() Authed = true, want false")
	}
	if c.Token != "" {
		t.Errorf("New() Token = %v, want empty string", c.Token)
	}
}

func TestIsSubscribed(t *testing.T) {
	c := client.New(&mockConn{})

	// Test when not subscribed
	if c.IsSubscribed("test") {
		t.Error("IsSubscribed() = true for non-subscribed topic")
	}

	// Test when subscribed
	c.Subscribe("test")
	if !c.IsSubscribed("test") {
		t.Error("IsSubscribed() = false for subscribed topic")
	}
}

func TestSubscribe(t *testing.T) {
	c := client.New(&mockConn{})

	// Test subscribing to a topic
	c.Subscribe("test")
	if !c.Subs["test"] {
		t.Error("Subscribe() did not add topic to subscriptions")
	}

	// Test subscribing to multiple topics
	c.Subscribe("test2")
	if !c.Subs["test2"] {
		t.Error("Subscribe() did not add second topic to subscriptions")
	}
}

func TestUnsubscribe(t *testing.T) {
	c := client.New(&mockConn{})

	// Subscribe to a topic first
	c.Subscribe("test")
	if !c.Subs["test"] {
		t.Error("Subscribe() did not add topic to subscriptions")
	}

	// Test unsubscribing
	c.Unsubscribe("test")
	if c.Subs["test"] {
		t.Error("Unsubscribe() did not remove topic from subscriptions")
	}

	// Test unsubscribing from non-existent topic
	c.Unsubscribe("nonexistent") // Should not panic
}

func TestUnsubscribeAll(t *testing.T) {
	c := client.New(&mockConn{})

	// Subscribe to multiple topics
	c.Subscribe("test1")
	c.Subscribe("test2")
	c.Subscribe("test3")

	// Test unsubscribing from all topics
	c.UnsubscribeAll()
	if len(c.Subs) != 0 {
		t.Errorf("UnsubscribeAll() left %d subscriptions, want 0", len(c.Subs))
	}
}

func TestWrite(t *testing.T) {
	conn := &mockConn{}
	c := client.New(conn)

	// Test writing a message
	message := "test message"
	err := c.Write(message)
	if err != nil {
		t.Errorf("Write() error = %v", err)
	}

	// Check if the message was written to the connection
	written := string(conn.writeData)
	if written != message {
		t.Errorf("Write() wrote %v, want %v", written, message)
	}
}

func TestClose(t *testing.T) {
	conn := &mockConn{}
	c := client.New(conn)

	// Test closing the connection
	err := c.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Check if the connection was closed
	if !conn.closed {
		t.Error("Close() did not close the connection")
	}
}

func TestConcurrentSubscribe(t *testing.T) {
	c := client.New(&mockConn{})
	done := make(chan bool)

	// Test concurrent subscriptions
	for i := 0; i < 100; i++ {
		go func(i int) {
			topic := "test" + string(rune(i))
			c.Subscribe(topic)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 100; i++ {
		<-done
	}

	// Verify all topics were subscribed
	for i := 0; i < 100; i++ {
		topic := "test" + string(rune(i))
		if !c.IsSubscribed(topic) {
			t.Errorf("Topic %s was not subscribed", topic)
		}
	}
}
