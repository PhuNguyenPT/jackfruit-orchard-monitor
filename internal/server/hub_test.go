package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	appConfig "GoApp/internal/config"

	"github.com/gorilla/websocket"
)

// Helper to create a base config for testing soil thresholds
func newTestConfig() *appConfig.Config {
	return &appConfig.Config{
		SoilDryValue: 3340,
		SoilWetValue: 1805,
	}
}

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

	h.BroadcastAirTempHumid("1", 30.1, 72.5)

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
	h.BroadcastSoilMoisture("0", 3340)

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

	h.BroadcastAirTempHumid("2", 28.0, 65.0)

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
	h.BroadcastAirTempHumid("1", 30.0, 70.0)
}
