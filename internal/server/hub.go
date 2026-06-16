package server

import (
	"bytes"
	"context"
	"log"
	"math"
	"strconv"
	"sync"
	"time"

	"GoApp/internal/database"
	"GoApp/internal/views"

	appConfig "GoApp/internal/config"
	"github.com/gorilla/websocket"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]string // conn -> lang
	cfg     *appConfig.Config
}

func NewHub(cfg *appConfig.Config) *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]string),
		cfg:     cfg,
	}
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

func (h *Hub) BroadcastAirTempHumid(addr string, temperature, humidity float32) {
	addrInt, err := strconv.ParseInt(addr, 10, 16)
	if err != nil {
		log.Printf("[Hub] invalid addr %q: %v", addr, err)
		return
	}

	row := database.GetLatestAirTempHumidReadingsRow{
		Addr:        int16(addrInt),
		Temperature: int16(math.Round(float64(temperature) * 10)),
		Humidity:    int16(math.Round(float64(humidity) * 10)),
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

func (h *Hub) BroadcastSoilMoisture(addr string, raw int) {
	addrInt, err := strconv.ParseInt(addr, 10, 16)
	if err != nil {
		log.Printf("[Hub] invalid soil addr %q: %v", addr, err)
		return
	}

	// Construct the database row format expected by Templ
	row := database.GetLatestSoilMoistureReadingsRow{
		SensorIdx: int16(addrInt),
		Raw:       int16(raw),
		CreatedAt: time.Now(),
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for c, lang := range h.clients {
		var buf bytes.Buffer

		// Render the SoilCardOOB, passing in h.cfg
		if err := views.SoilCardOOB(row, lang, h.cfg.SoilDryValue, h.cfg.SoilWetValue).Render(context.Background(), &buf); err != nil {
			log.Printf("[Hub] render error: %v", err)
			continue
		}

		if err := c.WriteMessage(websocket.TextMessage, buf.Bytes()); err != nil {
			log.Printf("[Hub] write error, dropping client: %v", err)
			go h.unregister(c)
		}
	}
}
