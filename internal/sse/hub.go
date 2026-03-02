package sse

import "sync"

// Hub manages SSE client connections per webhook ID.
type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[chan []byte]struct{}
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]map[chan []byte]struct{}),
	}
}

// Subscribe registers a new subscriber for the given webhookID.
// Returns a buffered data channel and an unsubscribe cleanup function.
func (h *Hub) Subscribe(webhookID string) (chan []byte, func()) {
	ch := make(chan []byte, 10)

	h.mu.Lock()
	if h.clients[webhookID] == nil {
		h.clients[webhookID] = make(map[chan []byte]struct{})
	}
	h.clients[webhookID][ch] = struct{}{}
	h.mu.Unlock()

	unsubscribe := func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		if subs, ok := h.clients[webhookID]; ok {
			delete(subs, ch)
			if len(subs) == 0 {
				delete(h.clients, webhookID)
			}
		}
		close(ch)
	}

	return ch, unsubscribe
}

// Publish sends data to all subscriber channels for the given webhookID.
// Uses non-blocking sends to avoid slow clients blocking the publisher.
func (h *Hub) Publish(webhookID string, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for ch := range h.clients[webhookID] {
		select {
		case ch <- data:
		default:
		}
	}
}
