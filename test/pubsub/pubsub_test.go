package pubsub_test

import (
	"net"
	"testing"
	"time"

	"redix/pkg/auth"
	"redix/pkg/client"
	"redix/pkg/pubsub"
)

// mockConn is a mock implementation of net.Conn for testing
type mockConn struct {
	writeData []byte
	closed    bool
}

func (m *mockConn) Read(b []byte) (n int, err error)   { return 0, nil }
func (m *mockConn) Write(b []byte) (n int, err error)  { m.writeData = b; return len(b), nil }
func (m *mockConn) Close() error                       { m.closed = true; return nil }
func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestNew(t *testing.T) {
	ps := pubsub.New()
	// Test that we can subscribe to a topic
	conn := &mockConn{}
	c := client.New(conn)
	ps.Subscribe("test", c)
	if !c.IsSubscribed("test") {
		t.Error("New() instance cannot subscribe to topics")
	}
}

func TestSubscribe(t *testing.T) {
	ps := pubsub.New()
	conn := &mockConn{}
	c := client.New(conn)

	// Test subscribing to a topic
	ps.Subscribe("test", c)
	if !c.IsSubscribed("test") {
		t.Error("Subscribe() did not add topic to client subscriptions")
	}

	// Test subscribing to multiple topics
	ps.Subscribe("test2", c)
	if !c.IsSubscribed("test2") {
		t.Error("Subscribe() did not add second topic to client subscriptions")
	}

	// Verify subscriber count
	count := ps.GetSubscriberCount("test")
	if count != 1 {
		t.Errorf("GetSubscriberCount() = %d, want 1", count)
	}
}

func TestUnsubscribe(t *testing.T) {
	ps := pubsub.New()
	conn := &mockConn{}
	c := client.New(conn)

	// Subscribe to a topic first
	ps.Subscribe("test", c)
	if !c.IsSubscribed("test") {
		t.Error("Subscribe() did not add topic to client subscriptions")
	}

	// Test unsubscribing
	ps.Unsubscribe("test", c)
	if c.IsSubscribed("test") {
		t.Error("Unsubscribe() did not remove topic from client subscriptions")
	}

	// Verify subscriber count
	count := ps.GetSubscriberCount("test")
	if count != 0 {
		t.Errorf("GetSubscriberCount() = %d, want 0", count)
	}
}

func TestPublish(t *testing.T) {
	ps := pubsub.New()
	conn1 := &mockConn{}
	conn2 := &mockConn{}
	c1 := client.New(conn1)
	c2 := client.New(conn2)

	// Set up test clients
	c1.Token = "token1"
	c2.Token = "token2"
	c1.Authed = true
	c2.Authed = true

	// Subscribe both clients to the same topic
	ps.Subscribe("test", c1)
	ps.Subscribe("test", c2)

	// Test publishing with token1
	count := ps.Publish("test", "hello", "token1")
	if count != 1 {
		t.Errorf("Publish() count = %d, want 1", count)
	}
	if string(conn1.writeData) == "" {
		t.Error("Publish() did not send message to token1 subscriber")
	}
	if string(conn2.writeData) != "" {
		t.Error("Publish() sent message to token2 subscriber")
	}

	// Test publishing with master token
	conn1.writeData = nil
	conn2.writeData = nil
	count = ps.Publish("test", "hello", auth.MasterToken)
	if count != 2 {
		t.Errorf("Publish() count = %d, want 2", count)
	}
	if string(conn1.writeData) == "" {
		t.Error("Publish() did not send message to token1 subscriber with master token")
	}
	if string(conn2.writeData) == "" {
		t.Error("Publish() did not send message to token2 subscriber with master token")
	}
}

func TestDisconnectToken(t *testing.T) {
	ps := pubsub.New()
	conn1 := &mockConn{}
	conn2 := &mockConn{}
	c1 := client.New(conn1)
	c2 := client.New(conn2)

	// Set up test clients
	c1.Token = "token1"
	c2.Token = "token1" // Same token
	c1.Authed = true
	c2.Authed = true

	// Subscribe both clients to different topics
	ps.Subscribe("test1", c1)
	ps.Subscribe("test2", c2)

	// Verify initial subscriber counts
	if ps.GetSubscriberCount("test1") != 1 {
		t.Error("Initial subscriber count for test1 is not 1")
	}
	if ps.GetSubscriberCount("test2") != 1 {
		t.Error("Initial subscriber count for test2 is not 1")
	}

	// Test disconnecting token1
	disconnected := ps.DisconnectToken("token1")
	if disconnected != 2 {
		t.Errorf("DisconnectToken() count = %d, want 2", disconnected)
	}
	if !conn1.closed {
		t.Error("DisconnectToken() did not close first client")
	}
	if !conn2.closed {
		t.Error("DisconnectToken() did not close second client")
	}

	// Verify subscriber counts after disconnect
	if ps.GetSubscriberCount("test1") != 0 {
		t.Error("Subscriber count for test1 is not 0 after disconnect")
	}
	if ps.GetSubscriberCount("test2") != 0 {
		t.Error("Subscriber count for test2 is not 0 after disconnect")
	}
}

func TestGetSubscriberCount(t *testing.T) {
	ps := pubsub.New()
	conn1 := &mockConn{}
	conn2 := &mockConn{}
	c1 := client.New(conn1)
	c2 := client.New(conn2)

	// Test empty topic
	count := ps.GetSubscriberCount("test")
	if count != 0 {
		t.Errorf("GetSubscriberCount() = %d, want 0", count)
	}

	// Subscribe clients
	ps.Subscribe("test", c1)
	ps.Subscribe("test", c2)

	// Test populated topic
	count = ps.GetSubscriberCount("test")
	if count != 2 {
		t.Errorf("GetSubscriberCount() = %d, want 2", count)
	}
}

func TestClientIsolation(t *testing.T) {
	ps := pubsub.New()
	conn1 := &mockConn{}
	conn2 := &mockConn{}
	conn3 := &mockConn{}
	c1 := client.New(conn1)
	c2 := client.New(conn2)
	c3 := client.New(conn3)

	// Set up test clients with different tokens
	c1.Token = "token1"
	c2.Token = "token2"
	c3.Token = "token1" // Same token as c1
	c1.Authed = true
	c2.Authed = true
	c3.Authed = true

	// Subscribe all clients to the same topic
	ps.Subscribe("test", c1)
	ps.Subscribe("test", c2)
	ps.Subscribe("test", c3)

	// Test 1: Publish with token1 - should only go to c1 and c3
	count := ps.Publish("test", "hello-token1", "token1")
	if count != 2 {
		t.Errorf("Publish() count = %d, want 2", count)
	}
	if string(conn1.writeData) == "" {
		t.Error("Publish() did not send message to token1 subscriber (c1)")
	}
	if string(conn2.writeData) != "" {
		t.Error("Publish() incorrectly sent message to token2 subscriber (c2)")
	}
	if string(conn3.writeData) == "" {
		t.Error("Publish() did not send message to token1 subscriber (c3)")
	}

	// Reset write data
	conn1.writeData = nil
	conn2.writeData = nil
	conn3.writeData = nil

	// Test 2: Publish with token2 - should only go to c2
	count = ps.Publish("test", "hello-token2", "token2")
	if count != 1 {
		t.Errorf("Publish() count = %d, want 1", count)
	}
	if string(conn1.writeData) != "" {
		t.Error("Publish() incorrectly sent message to token1 subscriber (c1)")
	}
	if string(conn2.writeData) == "" {
		t.Error("Publish() did not send message to token2 subscriber (c2)")
	}
	if string(conn3.writeData) != "" {
		t.Error("Publish() incorrectly sent message to token1 subscriber (c3)")
	}

	// Reset write data
	conn1.writeData = nil
	conn2.writeData = nil
	conn3.writeData = nil

	// Test 3: Publish with master token - should go to all clients
	count = ps.Publish("test", "hello-master", auth.MasterToken)
	if count != 3 {
		t.Errorf("Publish() count = %d, want 3", count)
	}
	if string(conn1.writeData) == "" {
		t.Error("Publish() did not send message to token1 subscriber (c1) with master token")
	}
	if string(conn2.writeData) == "" {
		t.Error("Publish() did not send message to token2 subscriber (c2) with master token")
	}
	if string(conn3.writeData) == "" {
		t.Error("Publish() did not send message to token1 subscriber (c3) with master token")
	}
}

func TestUnauthenticatedClientIsolation(t *testing.T) {
	ps := pubsub.New()
	conn1 := &mockConn{}
	conn2 := &mockConn{}
	c1 := client.New(conn1)
	c2 := client.New(conn2)

	// Set up test clients - one authenticated, one not
	c1.Token = "token1"
	c1.Authed = true
	c2.Authed = false // Unauthenticated client

	// Subscribe both clients to the same topic
	ps.Subscribe("test", c1)
	ps.Subscribe("test", c2)

	// Test 1: Publish with token1 - should only go to c1
	count := ps.Publish("test", "hello-token1", "token1")
	if count != 1 {
		t.Errorf("Publish() count = %d, want 1", count)
	}
	if string(conn1.writeData) == "" {
		t.Error("Publish() did not send message to authenticated subscriber")
	}
	if string(conn2.writeData) != "" {
		t.Error("Publish() incorrectly sent message to unauthenticated subscriber")
	}

	// Reset write data
	conn1.writeData = nil
	conn2.writeData = nil

	// Test 2: Publish with master token - should still only go to authenticated client
	count = ps.Publish("test", "hello-master", auth.MasterToken)
	if count != 1 {
		t.Errorf("Publish() count = %d, want 1", count)
	}
	if string(conn1.writeData) == "" {
		t.Error("Publish() did not send message to authenticated subscriber with master token")
	}
	if string(conn2.writeData) != "" {
		t.Error("Publish() incorrectly sent message to unauthenticated subscriber with master token")
	}
}

func TestMultiTopicIsolation(t *testing.T) {
	ps := pubsub.New()
	conn1 := &mockConn{}
	conn2 := &mockConn{}
	c1 := client.New(conn1)
	c2 := client.New(conn2)

	// Set up test clients with different tokens
	c1.Token = "token1"
	c2.Token = "token2"
	c1.Authed = true
	c2.Authed = true

	// Subscribe clients to different topics
	ps.Subscribe("topic1", c1)
	ps.Subscribe("topic2", c2)

	// Test 1: Publish to topic1 with token1
	count := ps.Publish("topic1", "hello-topic1", "token1")
	if count != 1 {
		t.Errorf("Publish() count = %d, want 1", count)
	}
	if string(conn1.writeData) == "" {
		t.Error("Publish() did not send message to topic1 subscriber")
	}
	if string(conn2.writeData) != "" {
		t.Error("Publish() incorrectly sent message to topic2 subscriber")
	}

	// Reset write data
	conn1.writeData = nil
	conn2.writeData = nil

	// Test 2: Publish to topic2 with token2
	count = ps.Publish("topic2", "hello-topic2", "token2")
	if count != 1 {
		t.Errorf("Publish() count = %d, want 1", count)
	}
	if string(conn1.writeData) != "" {
		t.Error("Publish() incorrectly sent message to topic1 subscriber")
	}
	if string(conn2.writeData) == "" {
		t.Error("Publish() did not send message to topic2 subscriber")
	}
}
