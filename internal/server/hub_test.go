package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// newTestConn spins up a server that upgrades and registers the *server-side*
// connection with h, then dials a client conn against it. BroadcastAirTempHumid sent to
// h will arrive on the returned client conn, mirroring production wiring.
func newTestConn(t *testing.T, h *Hub, lang string) (*websocket.Conn, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		up := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
		conn, err := up.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		h.register(conn, lang)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				h.unregister(conn)
				return
			}
		}
	}))
	url := "ws" + strings.TrimPrefix(srv.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	return conn, func() { conn.Close(); srv.Close() }
}

func TestHub_RegisterUnregister(t *testing.T) {
	cfg := newTestConfig()
	h := NewHub(cfg)
	conn, cleanup := newTestConn(t, h, "en")
	defer cleanup()

	waitForClientCount(t, h, 1)

	conn.Close()

	waitForClientCount(t, h, 0)
}

func TestHub_BroadcastAirTempHumid(t *testing.T) {
	cfg := newTestConfig()
	h := NewHub(cfg)
	conn, cleanup := newTestConn(t, h, "en")
	defer cleanup()

	waitForClientCount(t, h, 1)

	h.BroadcastAirTempHumid("1", 30.1, 72.5, time.Now())

	if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	html := string(msg)
	if !strings.Contains(html, `id="sensor-1"`) {
		t.Errorf("missing OOB target id, got: %s", html)
	}
	if !strings.Contains(html, `hx-swap-oob="true"`) {
		t.Errorf("missing hx-swap-oob attribute, got: %s", html)
	}
	if !strings.Contains(html, "30.1") || !strings.Contains(html, "72.5") {
		t.Errorf("missing reading values, got: %s", html)
	}
}

func TestHub_BroadcastSoilMoisture(t *testing.T) {
	cfg := newTestConfig() // dry: 3340, wet: 1805
	h := NewHub(cfg)
	conn, cleanup := newTestConn(t, h, "en")
	defer cleanup()

	waitForClientCount(t, h, 1)

	// Sending a raw value of 3340 should calculate out to exactly 0.0% moisture
	h.BroadcastSoilMoisture("0", 3340, time.Now())

	if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	html := string(msg)
	if !strings.Contains(html, `id="soil-0"`) {
		t.Errorf("missing OOB target id, got: %s", html)
	}
	if !strings.Contains(html, `hx-swap-oob="true"`) {
		t.Errorf("missing hx-swap-oob attribute, got: %s", html)
	}
	if !strings.Contains(html, "Raw: 3340") {
		t.Errorf("missing raw value label, got: %s", html)
	}
	if !strings.Contains(html, "0.0") {
		t.Errorf("missing expected calculated percentage (0.0%%), got: %s", html)
	}
}

func TestHub_BroadcastAirTempHumid_MultipleClients(t *testing.T) {
	cfg := newTestConfig()
	h := NewHub(cfg)

	conn1, cleanup1 := newTestConn(t, h, "en")
	defer cleanup1()
	conn2, cleanup2 := newTestConn(t, h, "en")
	defer cleanup2()

	waitForClientCount(t, h, 2)

	h.BroadcastAirTempHumid("2", 28.0, 65.0, time.Now())

	for _, conn := range []*websocket.Conn{conn1, conn2} {
		if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
			t.Fatalf("SetReadDeadline: %v", err)
		}
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("read: %v", err)
		}
		if !strings.Contains(string(msg), `id="sensor-2"`) {
			t.Errorf("missing OOB target id, got: %s", msg)
		}
	}
}

func TestHub_BroadcastAirTempHumid_DropsFailingClient(t *testing.T) {
	cfg := newTestConfig()
	h := NewHub(cfg)
	_, cleanup := newTestConn(t, h, "en")

	waitForClientCount(t, h, 1)

	cleanup() // close before broadcast

	waitForClientCount(t, h, 0)

	// Broadcasting with no clients should be a no-op, not a panic.
	h.BroadcastAirTempHumid("1", 30.0, 70.0, time.Now())
}

// ---------------------------------------------------------------------------
// Chart subscription helpers
// ---------------------------------------------------------------------------

func newSHT40ChartConn(t *testing.T, h *Hub, addr int16, lang string) (*websocket.Conn, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		up := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
		conn, err := up.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		h.registerSHT40Chart(conn, addr, lang)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				h.unregisterSHT40Chart(conn)
				return
			}
		}
	}))
	url := "ws" + strings.TrimPrefix(srv.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	return conn, func() { conn.Close(); srv.Close() }
}

func newSoilChartConn(t *testing.T, h *Hub, idx int16, lang string) (*websocket.Conn, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		up := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
		conn, err := up.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		h.registerSoilChart(conn, idx, lang)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				h.unregisterSoilChart(conn)
				return
			}
		}
	}))
	url := "ws" + strings.TrimPrefix(srv.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	return conn, func() { conn.Close(); srv.Close() }
}

// ---------------------------------------------------------------------------
// soilPct unit tests
// ---------------------------------------------------------------------------
func TestSoilPct(t *testing.T) {
	tests := []struct {
		raw  int16
		dry  int
		wet  int
		want float32
	}{
		{3340, 3340, 1805, 0},   // at dry — clamp 0
		{4000, 3340, 1805, 0},   // above dry — clamp 0
		{1805, 3340, 1805, 100}, // at wet — clamp 100
		{1000, 3340, 1805, 100}, // below wet — clamp 100
		{2550, 3300, 1800, 50},  // exact midpoint: (3300+1800)/2=2550, range=1500
		{3340, 3340, 3340, 0},   // equal dry/wet — guard divide-by-zero
	}
	for _, tt := range tests {
		got := soilPct(tt.raw, tt.dry, tt.wet)
		if got != tt.want {
			t.Errorf("soilPct(%d, %d, %d) = %.4f, want %.4f",
				tt.raw, tt.dry, tt.wet, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// SHT40 chart subscription lifecycle
// ---------------------------------------------------------------------------

func TestHub_SHT40Chart_RegisterAndUnregister(t *testing.T) {
	cfg := newTestConfig()
	h := NewHub(cfg)

	conn, cleanup := newSHT40ChartConn(t, h, 1, "en")

	waitForSHT40ChartCount(t, h, 1)
	conn.Close()
	waitForSHT40ChartCount(t, h, 0)
	cleanup()
}

func TestHub_SHT40Chart_ReceivesJSONPoint(t *testing.T) {
	cfg := newTestConfig()
	h := NewHub(cfg)
	conn, cleanup := newSHT40ChartConn(t, h, 1, "en")
	defer cleanup()

	waitForSHT40ChartCount(t, h, 1)

	h.BroadcastAirTempHumid("1", 32.5, 63.0, time.Now())

	if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	var pt map[string]any
	if err := json.Unmarshal(msg, &pt); err != nil {
		t.Fatalf("expected JSON chart point, got: %s", msg)
	}
	if _, ok := pt["t"]; !ok {
		t.Errorf("missing 't' field, got: %s", msg)
	}
	if _, ok := pt["temp"]; !ok {
		t.Errorf("missing 'temp' field, got: %s", msg)
	}
	if _, ok := pt["humid"]; !ok {
		t.Errorf("missing 'humid' field, got: %s", msg)
	}
}

func TestHub_SHT40Chart_FiltersWrongAddr(t *testing.T) {
	cfg := newTestConfig()
	h := NewHub(cfg)
	conn, cleanup := newSHT40ChartConn(t, h, 1, "en") // subscribed to addr 1
	defer cleanup()

	waitForSHT40ChartCount(t, h, 1)

	h.BroadcastAirTempHumid("2", 28.0, 55.0, time.Now()) // broadcast for addr 2

	if err := conn.SetReadDeadline(time.Now().Add(300 * time.Millisecond)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}
	_, _, err := conn.ReadMessage()
	if err == nil {
		t.Errorf("addr filter broken: subscriber for addr=1 received broadcast for addr=2")
	}
}

func TestHub_SHT40Chart_MultipleSubscribers(t *testing.T) {
	cfg := newTestConfig()
	h := NewHub(cfg)

	conn1, cleanup1 := newSHT40ChartConn(t, h, 1, "en")
	defer cleanup1()
	conn2, cleanup2 := newSHT40ChartConn(t, h, 1, "vi")
	defer cleanup2()

	waitForSHT40ChartCount(t, h, 2)

	h.BroadcastAirTempHumid("1", 31.0, 60.0, time.Now())

	for _, conn := range []*websocket.Conn{conn1, conn2} {
		if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
			t.Fatalf("SetReadDeadline: %v", err)
		}
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("read: %v", err)
		}
		var pt map[string]any
		if err := json.Unmarshal(msg, &pt); err != nil {
			t.Errorf("expected JSON, got: %s", msg)
		}
	}
}

func TestHub_SHT40Chart_DropsFailingConn(t *testing.T) {
	cfg := newTestConfig()
	h := NewHub(cfg)
	_, cleanup := newSHT40ChartConn(t, h, 1, "en")

	waitForSHT40ChartCount(t, h, 1)
	cleanup() // close before broadcast
	waitForSHT40ChartCount(t, h, 0)

	// Should be a no-op, not a panic.
	h.BroadcastAirTempHumid("1", 30.0, 60.0, time.Now())
}

// ---------------------------------------------------------------------------
// Soil chart subscription lifecycle
// ---------------------------------------------------------------------------

func TestHub_SoilChart_RegisterAndUnregister(t *testing.T) {
	cfg := newTestConfig()
	h := NewHub(cfg)

	conn, cleanup := newSoilChartConn(t, h, 0, "en")

	waitForSoilChartCount(t, h, 1)
	conn.Close()
	waitForSoilChartCount(t, h, 0)
	cleanup()
}

func TestHub_SoilChart_ReceivesJSONPoint(t *testing.T) {
	cfg := newTestConfig()
	h := NewHub(cfg)
	conn, cleanup := newSoilChartConn(t, h, 0, "en")
	defer cleanup()

	waitForSoilChartCount(t, h, 1)

	h.BroadcastSoilMoisture("0", 2572, time.Now())

	if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	var pt map[string]any
	if err := json.Unmarshal(msg, &pt); err != nil {
		t.Fatalf("expected JSON chart point, got: %s", msg)
	}
	if _, ok := pt["t"]; !ok {
		t.Errorf("missing 't' field, got: %s", msg)
	}
	if _, ok := pt["pct"]; !ok {
		t.Errorf("missing 'pct' field, got: %s", msg)
	}
}

func TestHub_SoilChart_FiltersWrongIdx(t *testing.T) {
	cfg := newTestConfig()
	h := NewHub(cfg)
	conn, cleanup := newSoilChartConn(t, h, 0, "en") // subscribed to idx 0
	defer cleanup()

	waitForSoilChartCount(t, h, 1)

	h.BroadcastSoilMoisture("1", 2000, time.Now()) // broadcast for idx 1

	if err := conn.SetReadDeadline(time.Now().Add(300 * time.Millisecond)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}
	_, _, err := conn.ReadMessage()
	if err == nil {
		t.Errorf("idx filter broken: subscriber for idx=0 received broadcast for idx=1")
	}
}

func TestHub_SoilChart_DropsFailingConn(t *testing.T) {
	cfg := newTestConfig()
	h := NewHub(cfg)
	_, cleanup := newSoilChartConn(t, h, 0, "en")

	waitForSoilChartCount(t, h, 1)
	cleanup() // close before broadcast
	waitForSoilChartCount(t, h, 0)

	// Should be a no-op, not a panic.
	h.BroadcastSoilMoisture("0", 2000, time.Now())
}

// ---------------------------------------------------------------------------
// Backfill (gap-fill on reconnect)
// ---------------------------------------------------------------------------
func TestHub_PushSHT40Backfill_NoRowsIsNoOp(t *testing.T) {
	cfg := newTestConfig()
	h := NewHub(cfg)
	conn, cleanup := newSHT40ChartConn(t, h, 1, "en")
	defer cleanup()

	waitForSHT40ChartCount(t, h, 1)

	db := &mockDB{}
	since := time.Now().Add(1 * time.Minute) // future timestamp — nothing should be "newer"

	h.pushSHT40Backfill(context.Background(), db, conn, 1, since, "en")

	if err := conn.SetReadDeadline(time.Now().Add(300 * time.Millisecond)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}
	_, _, err := conn.ReadMessage()
	if err == nil {
		t.Errorf("expected no message when backfill has zero rows, but got one")
	}
}
func TestHub_PushSoilBackfill_NoRowsIsNoOp(t *testing.T) {
	cfg := newTestConfig()
	h := NewHub(cfg)
	conn, cleanup := newSoilChartConn(t, h, 0, "en")
	defer cleanup()

	waitForSoilChartCount(t, h, 1)

	db := &mockDB{}
	since := time.Now().Add(1 * time.Minute)

	h.pushSoilBackfill(context.Background(), db, conn, 0, since, "en")

	if err := conn.SetReadDeadline(time.Now().Add(300 * time.Millisecond)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}
	_, _, err := conn.ReadMessage()
	if err == nil {
		t.Errorf("expected no message when backfill has zero rows, but got one")
	}
}
