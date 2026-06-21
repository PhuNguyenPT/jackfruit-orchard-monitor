package server

import (
	server "GoApp/internal/config"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var testHandler http.Handler

func newTestAppConfig() *server.Config {
	return &server.Config{
		AppEnv:   server.EnvTest,
		GinMode:  gin.TestMode,
		BaseURLs: []string{"http://localhost:8080"},
	}
}

func TestMain(m *testing.M) {
	if err := os.Chdir("../../"); err != nil {
		log.Fatalf("failed to change directory: %v", err)
	}
	s := &Server{
		db:  &mockDB{},
		cfg: newTestAppConfig(),
	}
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
