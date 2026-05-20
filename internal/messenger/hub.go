package messenger

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type Hub struct {
	clients    map[int64]*Client
	register   chan *Client
	unregister chan *Client
	direct     chan DirectMessage
	done       chan struct{}
}

type DirectMessage struct {
	From    int64  `json:"from"`
	To      int64  `json:"to"`
	Message string `json:"message"`
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		direct:     make(chan DirectMessage),
		clients:    make(map[int64]*Client),
		done:       make(chan struct{}),
	}
}

func (h *Hub) Stop() {
	close(h.done)
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			h.shutdown()
			return
		case <-h.done:
			h.shutdown()
			return
		case client := <-h.register:
			if oldClient, ok := h.clients[client.id]; ok {
				close(oldClient.send)
			}
			h.clients[client.id] = client
		case client := <-h.unregister:
			if c, ok := h.clients[client.id]; ok && c == client {
				delete(h.clients, client.id)
				close(client.send)
			}
		case msg := <-h.direct:
			if client, ok := h.clients[msg.To]; ok {
				msgBytes, err := json.Marshal(msg)
				if err != nil {
					hlog.Errorf("hub: marshal direct message: %v", err)
					continue
				}
				select {
				case client.send <- msgBytes:
				default:
					close(client.send)
					delete(h.clients, client.id)
				}
			}
		}
	}
}

func (h *Hub) shutdown() {
	for id, client := range h.clients {
		close(client.send)
		delete(h.clients, id)
	}
	hlog.Info("hub: shutdown complete")
}
