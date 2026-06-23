package server

import (
	"bytes"
	"context"
	"log"
	"maps"
	"math"
	"strconv"
	"sync"
	"time"

	"GoApp/internal/database"
	"GoApp/internal/views"

	config "GoApp/internal/config"
	"GoApp/internal/model"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]string    // conn -> lang
	devices map[string]model.DeviceStatus // clientID (MQTT username) → status
	cfg     *config.Config
}

func NewHub(cfg *config.Config) *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]string),
		devices: make(map[string]model.DeviceStatus),
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

func (h *Hub) BroadcastAirTempHumid(addr string, temperature, humidity float32, createdAt time.Time) {
	addrInt, err := strconv.ParseInt(addr, 10, 16)
	if err != nil {
		log.Printf("[Hub] invalid addr %q: %v", addr, err)
		return
	}

	row := database.GetLatestAirTempHumidReadingsRow{
		Addr:        int16(addrInt),
		Temperature: int16(math.Round(float64(temperature) * 10)),
		Humidity:    int16(math.Round(float64(humidity) * 10)),
		CreatedAt:   createdAt,
	}

	h.mu.RLock()
	defer h.mu.RUnlock()
	var buf bytes.Buffer
	for c, lang := range h.clients {
		buf.Reset()
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

func (h *Hub) BroadcastSoilMoisture(addr string, raw int, createdAt time.Time) {
	addrInt, err := strconv.ParseInt(addr, 10, 16)
	if err != nil {
		log.Printf("[Hub] invalid soil addr %q: %v", addr, err)
		return
	}

	// Construct the database row format expected by Templ
	row := database.GetLatestSoilMoistureReadingsRow{
		SensorIdx: int16(addrInt),
		Raw:       int16(raw),
		CreatedAt: createdAt,
	}

	h.mu.RLock()
	defer h.mu.RUnlock()
	var buf bytes.Buffer
	for c, lang := range h.clients {
		buf.Reset()
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

func (h *Hub) BroadcastDeviceStatus(clientID string, connected bool) {
	h.mu.Lock()
	h.devices[clientID] = model.DeviceStatus{Connected: connected, UpdatedAt: time.Now()}
	devicesCopy := make(map[string]model.DeviceStatus, len(h.devices))
	maps.Copy(devicesCopy, h.devices)
	h.mu.Unlock()

	h.mu.RLock()
	defer h.mu.RUnlock()
	var buf bytes.Buffer
	for c, lang := range h.clients {
		buf.Reset()
		if err := views.DeviceStatusListOOB(devicesCopy, lang).Render(context.Background(), &buf); err != nil {
			log.Printf("[Hub] render device status error: %v", err)
			continue
		}
		if err := c.WriteMessage(websocket.TextMessage, buf.Bytes()); err != nil {
			log.Printf("[Hub] write error, dropping client: %v", err)
			go h.unregister(c)
		}
	}
}

func (h *Hub) GetDeviceStatuses() map[string]model.DeviceStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make(map[string]model.DeviceStatus, len(h.devices))
	maps.Copy(out, h.devices)
	return out
}

func (h *Hub) pushDeviceStatusesToClient(conn *websocket.Conn, lang string) {
	h.mu.RLock()
	if len(h.devices) == 0 {
		h.mu.RUnlock()
		return
	}
	devicesCopy := make(map[string]model.DeviceStatus, len(h.devices))
	maps.Copy(devicesCopy, h.devices)
	h.mu.RUnlock()

	var buf bytes.Buffer
	if err := views.DeviceStatusListOOB(devicesCopy, lang).Render(context.Background(), &buf); err != nil {
		log.Printf("[Hub] render device status list error: %v", err)
		return
	}
	_ = conn.WriteMessage(websocket.TextMessage, buf.Bytes())
}
