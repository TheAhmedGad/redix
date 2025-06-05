package pubsub

import (
	"net"
	"sync"

	"redix/pkg/auth"
	"redix/pkg/client"
	"redix/pkg/protocol"
)

// PubSub handles the pub/sub functionality
type PubSub struct {
	subscribers map[string]map[net.Conn]*client.Client
	mu          sync.RWMutex
}

// New creates a new PubSub instance
func New() *PubSub {
	return &PubSub{
		subscribers: make(map[string]map[net.Conn]*client.Client),
	}
}

// Subscribe adds a client to a topic's subscribers
func (p *PubSub) Subscribe(topic string, c *client.Client) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.subscribers[topic] == nil {
		p.subscribers[topic] = make(map[net.Conn]*client.Client)
	}
	p.subscribers[topic][c.Conn] = c
	c.Subscribe(topic)
}

// Unsubscribe removes a client from a topic's subscribers
func (p *PubSub) Unsubscribe(topic string, c *client.Client) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if subs, ok := p.subscribers[topic]; ok {
		delete(subs, c.Conn)
		c.Unsubscribe(topic)
	}
}

// Publish sends a message to all subscribers of a topic
func (p *PubSub) Publish(topic, message string, publisherToken string) int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	subs := p.subscribers[topic]
	if subs == nil {
		return 0
	}

	count := 0
	for _, client := range subs {
		if client.Authed && (publisherToken == auth.MasterToken || client.Token == publisherToken) {
			client.Write(protocol.FormatMessage(topic, message))
			count++
		}
	}
	return count
}

// DisconnectToken disconnects all clients with a specific token
func (p *PubSub) DisconnectToken(targetToken string) int {
	p.mu.Lock()
	defer p.mu.Unlock()

	disconnectedConns := make(map[net.Conn]bool)
	disconnected := 0

	for _, subscribers := range p.subscribers {
		for conn, client := range subscribers {
			if client.Token == targetToken && !disconnectedConns[conn] {
				disconnectedConns[conn] = true
				client.Write(protocol.FormatError("disconnected by master"))
				client.Close()
				delete(subscribers, conn)
				disconnected++
			}
		}
	}

	return disconnected
}

// GetSubscriberCount returns the number of subscribers for a topic
func (p *PubSub) GetSubscriberCount(topic string) int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if subs, ok := p.subscribers[topic]; ok {
		return len(subs)
	}
	return 0
}
