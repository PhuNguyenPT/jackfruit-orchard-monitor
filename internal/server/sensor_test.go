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
	s := &Server{db: &mockDB{}, cfg: newTestConfig()}
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
	s := &Server{db: &erroringDB{&mockDB{}}, cfg: newTestConfig()}
	c, w := newSensorTestContext(http.MethodGet, "/sensors")

	s.sensorsPageHandler(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestSensorsGridHandler_WithReadings(t *testing.T) {
	s := &Server{db: &mockDB{}, cfg: newTestConfig()}
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
	s := &Server{db: &erroringDB{&mockDB{}}, cfg: newTestConfig()}
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

	s.hub.BroadcastAirTempHumid("1", 27.3, 55.5)

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
