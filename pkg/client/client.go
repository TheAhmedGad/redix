package client

import (
	"net"
	"sync"
)

// Client represents a connected client
type Client struct {
	Conn   net.Conn
	Token  string
	Authed bool
	Subs   map[string]bool
	mu     sync.RWMutex
}

// New creates a new client instance
func New(conn net.Conn) *Client {
	return &Client{
		Conn: conn,
		Subs: make(map[string]bool),
	}
}

// IsSubscribed checks if the client is subscribed to a topic
func (c *Client) IsSubscribed(topic string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Subs[topic]
}

// Subscribe adds a topic to the client's subscriptions
func (c *Client) Subscribe(topic string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Subs[topic] = true
}

// Unsubscribe removes a topic from the client's subscriptions
func (c *Client) Unsubscribe(topic string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.Subs, topic)
}

// UnsubscribeAll removes all topics from the client's subscriptions
func (c *Client) UnsubscribeAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Subs = make(map[string]bool)
}

// Write sends a message to the client
func (c *Client) Write(message string) error {
	_, err := c.Conn.Write([]byte(message))
	return err
}

// Close closes the client's connection
func (c *Client) Close() error {
	return c.Conn.Close()
}
