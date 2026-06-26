package server

import (
	config "GoApp/internal/config"
	"log/slog"
	"net/http"
	"strings"

	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var testHandler http.Handler

func newTestConfig() *config.Config {
	lv := &slog.LevelVar{}
	lv.Set(slog.LevelError)
	return &config.Config{
		AppEnv:       config.EnvTest,
		AppVersion:   "dev",
		BuildDate:    "2026-01-01",
		GinMode:      gin.TestMode,
		LogLevel:     lv,
		BaseURLs:     []string{"http://localhost:8080"},
		SoilDryValue: 3340,
		SoilWetValue: 1805,
	}
}
func newTestServer() *Server {
	cfg := newTestConfig()
	return &Server{
		db:  &mockDB{},
		cfg: cfg,
		hub: NewHub(cfg),
	}
}

func TestMain(m *testing.M) {
	if err := os.Chdir("../../"); err != nil {
		log.Fatalf("failed to change directory: %v", err)
	}
	s := newTestServer()
	testHandler = s.RegisterRoutes(s.cfg)
	os.Exit(m.Run())
}

func TestHomePageHandler(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	testHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}
}

func TestSitemapHandler(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	rr := httptest.NewRecorder()
	testHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}
	ct := rr.Header().Get("Content-Type")
	if ct != "application/xml; charset=utf-8" {
		t.Errorf("expected application/xml content-type, got %q", ct)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "2026-01-01") {
		t.Errorf("expected BuildDate in sitemap lastmod, got:\n%s", body)
	}
	if !strings.Contains(body, "http://localhost:8080/") {
		t.Errorf("expected base URL in sitemap, got:\n%s", body)
	}
}
