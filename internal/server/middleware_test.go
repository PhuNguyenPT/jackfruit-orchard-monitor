package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"GoApp/internal/rbac"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestRequirePermission_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := &mockDB{} // permissions is nil → every check returns false
	s := &Server{db: db}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", uuid.Must(uuid.NewV7()))
		c.Next()
	})
	router.GET("/admin/soil-calibration", s.requirePermission(rbac.ResourceSoilCalibration, rbac.ActionRead), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/soil-calibration", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestRequirePermission_Granted(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := &mockDB{permissions: map[string]bool{"soil_calibration:r": true}}
	s := &Server{db: db}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", uuid.Must(uuid.NewV7()))
		c.Next()
	})
	router.GET("/admin/soil-calibration", s.requirePermission(rbac.ResourceSoilCalibration, rbac.ActionRead), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/soil-calibration", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestRequirePermission_NoUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := &mockDB{}
	s := &Server{db: db}

	router := gin.New()
	// deliberately no userID set — simulates requirePermission used without authMiddleware ahead of it
	router.GET("/admin/soil-calibration", s.requirePermission(rbac.ResourceSoilCalibration, rbac.ActionRead), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/soil-calibration", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
