package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAboutPageHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/about", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	testHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}
	ct := rr.Header().Get("Content-Type")
	if ct != "text/html; charset=utf-8" {
		t.Errorf("expected text/html content-type, got %q", ct)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "About Prizm") {
		t.Errorf("expected About Prizm heading in body, got:\n%s", body)
	}
}
