package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"GoApp/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// erroringDB wraps mockDB but forces GetLatestAirTempHumidReadings to fail,
// so we can exercise the error branches without restating the mock.
type erroringDB struct {
	*mockDB
}

func (e *erroringDB) GetLatestAirTempHumidReadings(ctx context.Context) ([]database.GetLatestAirTempHumidReadingsRow, error) {
	return nil, errors.New("db down")
}

func (e *erroringDB) GetAirTempHumidReadingsByAddr(ctx context.Context, arg database.GetAirTempHumidReadingsByAddrParams) ([]database.GetAirTempHumidReadingsByAddrRow, error) {
	return nil, errors.New("db down")
}

func (e *erroringDB) GetSoilMoistureReadingsBySensorIdx(ctx context.Context, arg database.GetSoilMoistureReadingsBySensorIdxParams) ([]database.GetSoilMoistureReadingsBySensorIdxRow, error) {
	return nil, errors.New("db down")
}

func newSensorTestContext(method, path string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, nil)
	c.Set("lang", "en")
	c.Set("userName", "")
	return c, w
}

// ---------------------------------------------------------------------------
// sensorsPageHandler / sensorsGridHandler
// ---------------------------------------------------------------------------

func TestSensorsPageHandler_Success(t *testing.T) {
	s := newTestServer()
	c, w := newSensorTestContext(http.MethodGet, "/sensors")

	s.sensorsPageHandler(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	body := w.Body.String()
	if !strings.Contains(body, "<!doctype html>") {
		t.Errorf("expected full page document, got: %s", body)
	}
	for _, addr := range []string{"1", "2", "3"} {
		if !strings.Contains(body, `id="sensor-`+addr+`"`) {
			t.Errorf("missing card for sensor %s, got: %s", addr, body)
		}
	}
}

func TestSensorsPageHandler_DBError(t *testing.T) {
	s := newTestServer()
	s.db = &erroringDB{&mockDB{}}
	c, w := newSensorTestContext(http.MethodGet, "/sensors")

	s.sensorsPageHandler(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestSensorsGridHandler_WithReadings(t *testing.T) {
	s := newTestServer()
	c, w := newSensorTestContext(http.MethodGet, "/sensors/readings")

	s.sensorsGridHandler(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	body := w.Body.String()
	if !strings.Contains(body, "28.4") || !strings.Contains(body, "74.2") {
		t.Errorf("missing readings, got: %s", body)
	}
}

func TestSensorsGridHandler_DBError(t *testing.T) {
	s := newTestServer()
	s.db = &erroringDB{&mockDB{}}
	c, w := newSensorTestContext(http.MethodGet, "/sensors/readings")

	s.sensorsGridHandler(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// ---------------------------------------------------------------------------
// sensorsWSHandler — register/unregister + broadcast via the Hub
// ---------------------------------------------------------------------------

func waitForClientCount(t *testing.T, h *Hub, want int) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		h.mu.RLock()
		got := len(h.clients)
		h.mu.RUnlock()
		if got == want {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %d clients", want)
}

func newSensorWSServer() (*Server, *httptest.Server) {
	gin.SetMode(gin.TestMode)

	cfg := newTestConfig()
	s := &Server{db: &mockDB{}, cfg: cfg, hub: NewHub(cfg)}

	r := gin.New()
	r.GET("/sensors/ws", func(c *gin.Context) {
		c.Set("lang", "en")
		s.sensorsWSHandler(c)
	})
	return s, httptest.NewServer(r)
}

func TestSensorsWSHandler_RegisterAndUnregister(t *testing.T) {
	s, srv := newSensorWSServer()
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/sensors/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}

	waitForClientCount(t, s.hub, 1)
	conn.Close()
	waitForClientCount(t, s.hub, 0)
}

func TestSensorsWSHandler_BroadcastReachesClient(t *testing.T) {
	s, srv := newSensorWSServer()
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/sensors/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	waitForClientCount(t, s.hub, 1)

	s.hub.BroadcastAirTempHumid("1", 27.3, 55.5, time.Now())

	if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	html := string(msg)
	if !strings.Contains(html, `id="sensor-1"`) || !strings.Contains(html, `hx-swap-oob="true"`) {
		t.Errorf("unexpected OOB fragment: %s", html)
	}
	if !strings.Contains(html, "27.3") || !strings.Contains(html, "55.5") {
		t.Errorf("missing reading values, got: %s", html)
	}
}

const (
	deviceConnected    = true
	deviceDisconnected = false
)

func TestSensorsWSHandler_DeviceStatus_Connected(t *testing.T) {
	s, srv := newSensorWSServer()
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/sensors/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	waitForClientCount(t, s.hub, 1)

	s.hub.BroadcastDeviceStatus("esp32-nodemcu", deviceConnected)

	if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	html := string(msg)
	if !strings.Contains(html, `id="device-status-list"`) || !strings.Contains(html, `hx-swap-oob="true"`) {
		t.Errorf("expected full list container OOB swap, got: %s", html)
	}
	if !strings.Contains(html, "esp32-nodemcu") {
		t.Errorf("missing device clientID, got: %s", html)
	}
	if !strings.Contains(html, "Connected</span>") {
		t.Errorf("expected Connected status, got: %s", html)
	}
	if strings.Contains(html, "Disconnected</span>") {
		t.Errorf("did not expect Disconnected status, got: %s", html)
	}
}

func TestSensorsWSHandler_DeviceStatus_Disconnected(t *testing.T) {
	s, srv := newSensorWSServer()
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/sensors/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	waitForClientCount(t, s.hub, 1)

	s.hub.BroadcastDeviceStatus("esp32-nodemcu", deviceDisconnected)

	if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	html := string(msg)
	if !strings.Contains(html, `id="device-status-list"`) || !strings.Contains(html, `hx-swap-oob="true"`) {
		t.Errorf("expected full list container OOB swap, got: %s", html)
	}
	if !strings.Contains(html, "esp32-nodemcu") {
		t.Errorf("missing device clientID, got: %s", html)
	}
	if !strings.Contains(html, "Disconnected</span>") {
		t.Errorf("expected Disconnected status, got: %s", html)
	}
}

func TestSensorsWSHandler_DeviceStatus_FullCycle(t *testing.T) {
	s, srv := newSensorWSServer()
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/sensors/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	waitForClientCount(t, s.hub, 1)

	const clientID = "esp32-nodemcu"

	// The lifecycle under test: disconnected -> connected -> disconnected -> connected.
	transitions := []struct {
		name      string
		connected bool
	}{
		{"disconnected_to_connected", deviceConnected},
		{"connected_to_disconnected", deviceDisconnected},
		{"disconnected_to_connected_again", deviceConnected},
	}

	for _, tr := range transitions {
		t.Run(tr.name, func(t *testing.T) {
			s.hub.BroadcastDeviceStatus(clientID, tr.connected)

			if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
				t.Fatalf("SetReadDeadline: %v", err)
			}
			_, msg, err := conn.ReadMessage()
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			html := string(msg)

			if !strings.Contains(html, `id="device-status-list"`) || !strings.Contains(html, `hx-swap-oob="true"`) {
				t.Errorf("expected full list container OOB swap, got: %s", html)
			}
			if !strings.Contains(html, clientID) {
				t.Errorf("missing device clientID, got: %s", html)
			}

			wantTag, dontWantTag := "Disconnected</span>", "Connected</span>"
			if tr.connected {
				wantTag, dontWantTag = "Connected</span>", "Disconnected</span>"
			}
			if !strings.Contains(html, wantTag) {
				t.Errorf("expected %s, got: %s", wantTag, html)
			}
			if strings.Contains(html, dontWantTag) {
				t.Errorf("did not expect %s, got: %s", dontWantTag, html)
			}
		})
	}

	// Verify the Hub's internal state matches the last transition.
	s.hub.mu.RLock()
	status, ok := s.hub.devices[clientID]
	s.hub.mu.RUnlock()
	if !ok {
		t.Fatalf("expected device to be tracked in hub.devices")
	}
	want := transitions[len(transitions)-1].connected
	if status.Connected != want {
		t.Errorf("expected hub.devices to reflect Connected=%v after final transition, got %v", want, status.Connected)
	}
}

// ---------------------------------------------------------------------------
// SHT40HistoryPage / sht40HistoryHandler
// ---------------------------------------------------------------------------

func TestSHT40HistoryHandler_Success(t *testing.T) {
	s := newTestServer()
	c, w := newSensorTestContext(http.MethodGet, "/sensors/sht40/1")

	// Inject the path parameter expected by the handler
	c.Params = []gin.Param{{Key: "addr", Value: "1"}}

	s.sht40HistoryHandler(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.String()

	// Assert the script bundle and data container exist
	if !strings.Contains(body, `src="/public/sht40-history.min.js"`) {
		t.Errorf("missing script entry point for sht40 history chart")
	}
	if !strings.Contains(body, `id="chart-data"`) {
		t.Errorf("missing chart data element carrier")
	}

	// Verify the newly introduced localization attributes are present
	if !strings.Contains(body, "data-label-temp=") || !strings.Contains(body, "data-label-humid=") {
		t.Errorf("missing localized data attributes for SHT40 charts, got: %s", body)
	}
}

func TestSHT40HistoryHandler_DBError(t *testing.T) {
	s := newTestServer()
	s.db = &erroringDB{&mockDB{}}
	c, w := newSensorTestContext(http.MethodGet, "/sensors/sht40/1")
	c.Params = []gin.Param{{Key: "addr", Value: "1"}}

	s.sht40HistoryHandler(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// ---------------------------------------------------------------------------
// SoilHistoryPage / soilHistoryHandler
// ---------------------------------------------------------------------------

func TestSoilHistoryHandler_Success(t *testing.T) {
	s := newTestServer()
	c, w := newSensorTestContext(http.MethodGet, "/sensors/soil/0")

	// Inject the path parameter expected by the handler
	c.Params = []gin.Param{{Key: "idx", Value: "0"}}

	s.soilHistoryHandler(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.String()

	// Assert the script bundle and data container exist
	if !strings.Contains(body, `src="/public/soil-history.min.js"`) {
		t.Errorf("missing script entry point for soil history chart")
	}

	// Verify the newly introduced soil localization attribute is present
	if !strings.Contains(body, "data-label-soil=") {
		t.Errorf("missing localized data-label-soil attribute for soil chart, got: %s", body)
	}
}

func TestSoilHistoryHandler_DBError(t *testing.T) {
	s := newTestServer()
	s.db = &erroringDB{&mockDB{}}
	c, w := newSensorTestContext(http.MethodGet, "/sensors/soil/0")
	c.Params = []gin.Param{{Key: "idx", Value: "0"}}

	s.soilHistoryHandler(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
