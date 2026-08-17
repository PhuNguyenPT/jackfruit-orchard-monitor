package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	appConfig "GoApp/internal/config"

	"github.com/gin-gonic/gin"
)

func newTestSoilCalServer() *Server {
	cfg := &appConfig.Config{}
	return &Server{
		db:  &mockDB{},
		cfg: cfg,
		hub: NewHub(cfg),
	}
}

func TestSoilCalibrationUpdateHandler_InvalidSensorIdx(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := newTestSoilCalServer()

	router := gin.New()
	router.POST("/admin/soil-calibration", s.soilCalibrationUpdateHandler)

	body := strings.NewReader(url.Values{
		"sensor_idx": {"not-a-number"},
		"dry_value":  {"3000"},
		"wet_value":  {"1500"},
	}.Encode())
	req := httptest.NewRequest(http.MethodPost, "/admin/soil-calibration", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/admin/soil-calibration?status=error" {
		t.Errorf("expected error redirect, got %q", loc)
	}
}

func TestSoilCalibrationUpdateHandler_DryEqualsWet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := newTestSoilCalServer()

	router := gin.New()
	router.POST("/admin/soil-calibration", s.soilCalibrationUpdateHandler)

	body := strings.NewReader(url.Values{
		"sensor_idx": {"0"},
		"dry_value":  {"2000"},
		"wet_value":  {"2000"},
	}.Encode())
	req := httptest.NewRequest(http.MethodPost, "/admin/soil-calibration", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/admin/soil-calibration?status=error" {
		t.Errorf("expected error redirect, got %q", loc)
	}
}

func TestSoilCalibrationUpdateHandler_OutOfRange(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := newTestSoilCalServer()

	router := gin.New()
	router.POST("/admin/soil-calibration", s.soilCalibrationUpdateHandler)

	body := strings.NewReader(url.Values{
		"sensor_idx": {"0"},
		"dry_value":  {"5000"}, // > 4095
		"wet_value":  {"1500"},
	}.Encode())
	req := httptest.NewRequest(http.MethodPost, "/admin/soil-calibration", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/admin/soil-calibration?status=error" {
		t.Errorf("expected error redirect, got %q", loc)
	}
}

func TestSoilCalibrationUpdateHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := newTestSoilCalServer()

	router := gin.New()
	router.POST("/admin/soil-calibration", s.soilCalibrationUpdateHandler)

	body := strings.NewReader(url.Values{
		"sensor_idx": {"2"},
		"dry_value":  {"3000"},
		"wet_value":  {"1200"},
	}.Encode())
	req := httptest.NewRequest(http.MethodPost, "/admin/soil-calibration", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 redirect, got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/admin/soil-calibration?status=saved" {
		t.Errorf("expected success redirect with status=saved, got %q", loc)
	}

	// confirm the write hit the DB mock
	db := s.db.(*mockDB)
	saved, ok := db.soilCalibrations[2]
	if !ok {
		t.Fatal("expected sensor_idx 2 to be saved")
	}
	if saved.DryValue != 3000 || saved.WetValue != 1200 {
		t.Errorf("expected dry=3000 wet=1200, got dry=%d wet=%d", saved.DryValue, saved.WetValue)
	}

	// confirm the hub cache picked it up immediately (live WS/history reads use this, not the DB)
	cal := s.hub.GetSoilCalibration(2)
	if cal.Dry != 3000 || cal.Wet != 1200 {
		t.Errorf("expected hub cache dry=3000 wet=1200, got dry=%d wet=%d", cal.Dry, cal.Wet)
	}
}

func TestSoilCalibrationAdminPageHandler_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := newTestSoilCalServer()

	router := gin.New()
	router.GET("/admin/soil-calibration", s.soilCalibrationAdminPageHandler)

	req := httptest.NewRequest(http.MethodGet, "/admin/soil-calibration", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("expected html content type, got %q", ct)
	}
}

func TestHub_GetSoilCalibration_DefaultsFromConfig(t *testing.T) {
	cfg := &appConfig.Config{
		SoilDryValue: 3300,
		SoilWetValue: 1400,
	}
	hub := NewHub(cfg)

	// sensor 99 has never been calibrated — no DB row, no SetSoilCalibration call
	cal := hub.GetSoilCalibration(99)

	if cal.Dry != 3300 || cal.Wet != 1400 {
		t.Errorf("expected fallback to config defaults dry=3300 wet=1400, got dry=%d wet=%d", cal.Dry, cal.Wet)
	}
}

func TestHub_GetSoilCalibration_SavedOverridesConfig(t *testing.T) {
	cfg := &appConfig.Config{
		SoilDryValue: 3300,
		SoilWetValue: 1400,
	}
	hub := NewHub(cfg)
	hub.SetSoilCalibration(5, 2900, 1100)

	cal := hub.GetSoilCalibration(5)

	if cal.Dry != 2900 || cal.Wet != 1100 {
		t.Errorf("expected saved calibration to override config default, got dry=%d wet=%d", cal.Dry, cal.Wet)
	}
}

func TestSoilCalibrationUpdateHandler_HTMX_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := newTestSoilCalServer()

	router := gin.New()
	router.POST("/admin/soil-calibration", s.soilCalibrationUpdateHandler)

	body := strings.NewReader(url.Values{
		"sensor_idx": {"2"},
		"dry_value":  {"3000"},
		"wet_value":  {"1200"},
	}.Encode())
	req := httptest.NewRequest(http.MethodPost, "/admin/soil-calibration", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// htmx path renders the row inline instead of redirecting
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for htmx request, got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "" {
		t.Errorf("expected no redirect on htmx path, got Location %q", loc)
	}
	if !strings.Contains(w.Body.String(), `id="soil-row-2"`) {
		t.Errorf("expected row fragment for sensor 2, got body: %s", w.Body.String())
	}

	db := s.db.(*mockDB)
	saved, ok := db.soilCalibrations[2]
	if !ok || saved.DryValue != 3000 || saved.WetValue != 1200 {
		t.Errorf("expected sensor 2 saved dry=3000 wet=1200, got %+v (ok=%v)", saved, ok)
	}
}

func TestSoilCalibrationUpdateHandler_HTMX_ValidationError_RevertsRow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := newTestSoilCalServer()
	s.hub.SetSoilCalibration(2, 2800, 1300) // pre-existing known-good values

	router := gin.New()
	router.POST("/admin/soil-calibration", s.soilCalibrationUpdateHandler)

	body := strings.NewReader(url.Values{
		"sensor_idx": {"2"},
		"dry_value":  {"5000"}, // out of range
		"wet_value":  {"1300"},
	}.Encode())
	req := httptest.NewRequest(http.MethodPost, "/admin/soil-calibration", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	// row should show the last known-good values, not the rejected 5000
	if !strings.Contains(w.Body.String(), `value="2800"`) {
		t.Errorf("expected reverted dry value 2800 in fragment, got: %s", w.Body.String())
	}
}
