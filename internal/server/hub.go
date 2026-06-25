package server

import (
	"bytes"
	"context"
	"encoding/json"
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

var hubVietnamTZ = time.FixedZone("Asia/Ho_Chi_Minh", 7*60*60)

type chartSub struct {
	id   int16
	lang string
}
type Hub struct {
	mu          sync.RWMutex
	clients     map[*websocket.Conn]string    // conn -> lang
	devices     map[string]model.DeviceStatus // clientID (MQTT username) → status
	sht40Charts map[*websocket.Conn]chartSub  // conn -> addr filter
	soilCharts  map[*websocket.Conn]chartSub  // conn -> sensorIdx filter
	cfg         *config.Config
}

func NewHub(cfg *config.Config) *Hub {
	return &Hub{
		clients:     make(map[*websocket.Conn]string),
		devices:     make(map[string]model.DeviceStatus),
		sht40Charts: make(map[*websocket.Conn]chartSub),
		soilCharts:  make(map[*websocket.Conn]chartSub),
		cfg:         cfg,
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

func (h *Hub) registerSHT40Chart(conn *websocket.Conn, addr int16, lang string) {
	h.mu.Lock()
	h.sht40Charts[conn] = chartSub{id: addr, lang: lang}
	h.mu.Unlock()
}

func (h *Hub) unregisterSHT40Chart(conn *websocket.Conn) {
	h.mu.Lock()
	delete(h.sht40Charts, conn)
	h.mu.Unlock()
	conn.Close()
}

func (h *Hub) registerSoilChart(conn *websocket.Conn, idx int16, lang string) {
	h.mu.Lock()
	h.soilCharts[conn] = chartSub{id: idx, lang: lang}
	h.mu.Unlock()
}

func (h *Hub) unregisterSoilChart(conn *websocket.Conn) {
	h.mu.Lock()
	delete(h.soilCharts, conn)
	h.mu.Unlock()
	conn.Close()
}

func soilPct(raw int16, dryVal, wetVal int) float32 {
	if dryVal == wetVal {
		return 0
	}
	pct := float32(dryVal-int(raw)) / float32(dryVal-wetVal) * 100
	if pct < 0 {
		pct = 0
	} else if pct > 100 {
		pct = 100
	}
	return pct
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
	// --- live chart push ---
	for conn, sub := range h.sht40Charts {
		if sub.id != row.Addr {
			continue
		}
		format := "02-01 15:04:05"
		if sub.lang != "vi" {
			format = "01-02 15:04:05"
		}
		msg, err := json.Marshal(sht40Point{
			T:     createdAt.In(hubVietnamTZ).Format(format),
			Temp:  temperature,
			Humid: humidity,
		})
		if err != nil {
			log.Printf("[Hub] sht40 chart marshal error: %v", err)
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Printf("[Hub] sht40 chart write error, dropping conn: %v", err)
			go h.unregisterSHT40Chart(conn)
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
	// --- live chart push ---
	for conn, sub := range h.soilCharts {
		if sub.id != row.SensorIdx {
			continue
		}
		format := "02-01 15:04:05"
		if sub.lang != "vi" {
			format = "01-02 15:04:05"
		}
		msg, err := json.Marshal(soilPoint{
			T:   createdAt.In(hubVietnamTZ).Format(format),
			Pct: soilPct(int16(raw), h.cfg.SoilDryValue, h.cfg.SoilWetValue),
		})
		if err != nil {
			log.Printf("[Hub] soil chart marshal error: %v", err)
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Printf("[Hub] soil chart write error, dropping conn: %v", err)
			go h.unregisterSoilChart(conn)
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

type sht40Point struct {
	T     string  `json:"t"`
	Temp  float32 `json:"temp"`
	Humid float32 `json:"humid"`
}

type soilPoint struct {
	T   string  `json:"t"`
	Pct float32 `json:"pct"`
}

type chartBatch[T any] struct {
	Batch  bool `json:"batch"`
	Points []T  `json:"points"`
}

func (h *Hub) pushSHT40Backfill(ctx context.Context, db DB, conn *websocket.Conn, addr int16, since time.Time, lang string) {
	rows, err := db.GetAirTempHumidReadingsByAddrSince(ctx, database.GetAirTempHumidReadingsByAddrSinceParams{
		Addr:      addr,
		CreatedAt: since,
	})
	if err != nil {
		log.Printf("[Hub] sht40 backfill query error: %v", err)
		return
	}
	if len(rows) == 0 {
		return
	}

	format := "02-01 15:04:05"
	if lang != "vi" {
		format = "01-02 15:04:05"
	}

	points := make([]sht40Point, 0, len(rows))
	for _, r := range rows {
		points = append(points, sht40Point{
			T:     r.CreatedAt.In(hubVietnamTZ).Format(format),
			Temp:  float32(r.Temperature) / 10,
			Humid: float32(r.Humidity) / 10,
		})
	}

	msg, err := json.Marshal(chartBatch[sht40Point]{Batch: true, Points: points})
	if err != nil {
		log.Printf("[Hub] sht40 backfill marshal error: %v", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Printf("[Hub] sht40 backfill write error: %v", err)
	}
}

func (h *Hub) pushSoilBackfill(ctx context.Context, db DB, conn *websocket.Conn, idx int16, since time.Time, lang string) {
	rows, err := db.GetSoilMoistureReadingsBySensorIdxSince(ctx, database.GetSoilMoistureReadingsBySensorIdxSinceParams{
		SensorIdx: idx,
		CreatedAt: since,
	})
	if err != nil {
		log.Printf("[Hub] soil backfill query error: %v", err)
		return
	}
	if len(rows) == 0 {
		return
	}

	format := "02-01 15:04:05"
	if lang != "vi" {
		format = "01-02 15:04:05"
	}

	points := make([]soilPoint, 0, len(rows))
	for _, r := range rows {
		points = append(points, soilPoint{
			T:   r.CreatedAt.In(hubVietnamTZ).Format(format),
			Pct: soilPct(r.Raw, h.cfg.SoilDryValue, h.cfg.SoilWetValue),
		})
	}

	msg, err := json.Marshal(chartBatch[soilPoint]{Batch: true, Points: points})
	if err != nil {
		log.Printf("[Hub] soil backfill marshal error: %v", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Printf("[Hub] soil backfill write error: %v", err)
	}
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
