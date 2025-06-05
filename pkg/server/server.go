package server

import (
	"database/sql"
	"net"
	"strings"

	"redix/pkg/auth"
	"redix/pkg/client"
	"redix/pkg/protocol"
	"redix/pkg/pubsub"
)

// Server represents the main server instance
type Server struct {
	db      *sql.DB
	auth    *auth.Validator
	pubsub  *pubsub.PubSub
	handler *Handler
}

// New creates a new server instance
func New(db *sql.DB) *Server {
	validator := auth.NewValidator(db)
	ps := pubsub.New()
	handler := NewHandler(validator, ps)

	return &Server{
		db:      db,
		auth:    validator,
		pubsub:  ps,
		handler: handler,
	}
}

// Listen starts the server on the specified address
func (s *Server) Listen(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		c := client.New(conn)
		go s.handler.Handle(c)
	}
}

// Handler handles client connections and commands
type Handler struct {
	auth   *auth.Validator
	pubsub *pubsub.PubSub
}

// NewHandler creates a new command handler
func NewHandler(validator *auth.Validator, ps *pubsub.PubSub) *Handler {
	return &Handler{
		auth:   validator,
		pubsub: ps,
	}
}

// Handle processes client commands
func (h *Handler) Handle(c *client.Client) {
	defer c.Close()
	buf := make([]byte, 4096)

	for {
		n, err := c.Conn.Read(buf)
		if err != nil {
			break
		}

		cmd := protocol.ParseRESP(buf[:n])
		if len(cmd) == 0 {
			continue
		}

		switch strings.ToUpper(cmd[0]) {
		case "AUTH":
			if len(cmd) < 2 {
				c.Write(protocol.FormatError("wrong number of arguments for AUTH"))
				continue
			}
			token := cmd[1]
			if h.auth.IsValidToken(token) {
				c.Authed = true
				c.Token = token
				c.Write(protocol.FormatOK())
			} else {
				c.Write(protocol.FormatError("invalid token"))
			}

		case "DISCONNECT":
			if !c.Authed {
				c.Write(protocol.FormatNoAuth())
				continue
			}

			if !auth.IsMasterToken(c.Token) {
				c.Write(protocol.FormatError("only master token can disconnect clients"))
				continue
			}

			if len(cmd) < 2 {
				c.Write(protocol.FormatError("wrong number of arguments for DISCONNECT"))
				continue
			}

			disconnected := h.pubsub.DisconnectToken(cmd[1])
			c.Write(protocol.FormatInteger(disconnected))

		case "SUBSCRIBE":
			if !c.Authed {
				c.Write(protocol.FormatNoAuth())
				continue
			}

			for _, topic := range cmd[1:] {
				h.pubsub.Subscribe(topic, c)
				c.Write(protocol.FormatSubscribe(topic))
			}

		case "PUBLISH":
			if !c.Authed {
				c.Write(protocol.FormatNoAuth())
				continue
			}

			if len(cmd) < 3 {
				c.Write(protocol.FormatError("wrong number of arguments"))
				continue
			}

			topic, msg := cmd[1], cmd[2]
			count := h.pubsub.Publish(topic, msg, c.Token)
			c.Write(protocol.FormatInteger(count))

		default:
			c.Write(protocol.FormatError("unknown command"))
		}
	}
}
