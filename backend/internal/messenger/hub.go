package messenger

import (
	"encoding/json"

	"github.com/gliedabrennung/messenger-core/internal/usecase"
)

type Hub struct {
	clients    map[int64]*Client
	register   chan *Client
	unregister chan *Client
	direct     chan DirectMessage
	msgUseCase *usecase.MessageUseCase
}

type DirectMessage struct {
	Action  string `json:"action,omitempty"`
	From    int64  `json:"from"`
	To      int64  `json:"to"`
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
	Status  string `json:"status,omitempty"`
}

var hub *Hub

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		direct:     make(chan DirectMessage),
		clients:    make(map[int64]*Client),
	}
}

func InitHub(useCase *usecase.MessageUseCase) {
	hub = NewHub()
	hub.msgUseCase = useCase
	go hub.Run()
}

func (h *Hub) Run() {
	for {
		select {
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
