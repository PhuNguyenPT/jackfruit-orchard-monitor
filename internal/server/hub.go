package server

import (
	"bytes"
	"context"
	"log"
	"sync"
	"time"

	"GoApp/internal/database"
	"GoApp/internal/views"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]string // conn -> lang
}

func NewHub() *Hub {
	return &Hub{clients: make(map[*websocket.Conn]string)}
}

func (h *Hub) register(c *websocket.Conn, lang string) {
	h.mu.Lock()
	h.clients[c] = lang
	h.mu.Unlock()
}

func (h *Hub) unregister(c *websocket.Conn) {
	h.mu.Lock()
	delete(h.clients, c)
	h.mu.Unlock()
	c.Close()
}

func (h *Hub) Broadcast(addr string, temperature, humidity float32) {
	row := database.GetLatestReadingsRow{
		Addr:        addr,
		Temperature: temperature,
		Humidity:    humidity,
		CreatedAt:   time.Now(),
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for c, lang := range h.clients {
		var buf bytes.Buffer
		if err := views.SensorCardOOB(row, lang).Render(context.Background(), &buf); err != nil {
			log.Printf("[Hub] render error: %v", err)
			continue
		}
		if err := c.WriteMessage(websocket.TextMessage, buf.Bytes()); err != nil {
			log.Printf("[Hub] write error, dropping client: %v", err)
			go h.unregister(c)
		}
	}
}
