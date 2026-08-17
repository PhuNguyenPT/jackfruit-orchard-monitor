package server

import (
	"GoApp/internal/database"
	"GoApp/internal/views"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) soilCalibrationAdminPageHandler(c *gin.Context) {
	lang := getLangStr(c)

	readings, err := s.db.GetLatestSoilMoistureReadings(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	cals := make([]database.SoilCalibration, 0, len(readings))
	for _, r := range readings {
		cal := s.hub.GetSoilCalibration(r.SensorIdx)
		cals = append(cals, database.SoilCalibration{
			SensorIdx: r.SensorIdx,
			DryValue:  int16(cal.Dry),
			WetValue:  int16(cal.Wet),
		})
	}
	sort.Slice(cals, func(i, j int) bool { return cals[i].SensorIdx < cals[j].SensorIdx })

	status := c.Query("status") // "saved" | "error" | "" — no-JS fallback only

	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := views.SoilCalibrationAdminPage(cals, lang, getUserName(c), s.siteConfig(c), status).
		Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (s *Server) soilCalibrationUpdateHandler(c *gin.Context) {
	isHX := c.GetHeader("HX-Request") == "true"
	lang := getLangStr(c)

	idx, err := strconv.ParseInt(c.PostForm("sensor_idx"), 10, 16)
	if err != nil {
		// Can't identify which row to swap — sensor_idx is a hidden field we
		// control, so this only happens on a tampered/malformed request.
		if isHX {
			s.renderSoilCalToast(c, lang, "error")
			return
		}
		c.Redirect(http.StatusSeeOther, "/admin/soil-calibration?status=error")
		return
	}

	dry, err1 := strconv.ParseInt(c.PostForm("dry_value"), 10, 16)
	wet, err2 := strconv.ParseInt(c.PostForm("wet_value"), 10, 16)
	if err1 != nil || err2 != nil ||
		dry < 0 || dry > 4095 || wet < 0 || wet > 4095 || dry == wet {
		if isHX {
			s.renderSoilCalRow(c, int16(idx), lang, "error")
			return
		}
		c.Redirect(http.StatusSeeOther, "/admin/soil-calibration?status=error")
		return
	}

	saved, err := s.db.UpsertSoilCalibration(c.Request.Context(), database.UpsertSoilCalibrationParams{
		SensorIdx: int16(idx),
		DryValue:  int16(dry),
		WetValue:  int16(wet),
	})
	if err != nil {
		if isHX {
			s.renderSoilCalRow(c, int16(idx), lang, "error")
			return
		}
		c.Redirect(http.StatusSeeOther, "/admin/soil-calibration?status=error")
		return
	}

	s.hub.SetSoilCalibration(saved.SensorIdx, int(saved.DryValue), int(saved.WetValue))

	if isHX {
		c.Status(http.StatusOK)
		c.Header("Content-Type", "text/html; charset=utf-8")
		cal := database.SoilCalibration{
			SensorIdx: saved.SensorIdx,
			DryValue:  saved.DryValue,
			WetValue:  saved.WetValue,
		}
		if err := views.SoilCalibrationRow(cal, lang, "saved").Render(c.Request.Context(), c.Writer); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/soil-calibration?status=saved")
}

// renderSoilCalRow re-renders a single row from the hub cache (the same
// source of truth soilCalibrationAdminPageHandler uses), tagged with a
// status badge. Used on the htmx path when a write fails validation or the
// DB write itself errors, so the row reverts to the last known-good values
// rather than echoing back rejected input.
func (s *Server) renderSoilCalRow(c *gin.Context, sensorIdx int16, lang, rowStatus string) {
	calVal := s.hub.GetSoilCalibration(sensorIdx)
	cal := database.SoilCalibration{
		SensorIdx: sensorIdx,
		DryValue:  int16(calVal.Dry),
		WetValue:  int16(calVal.Wet),
	}
	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := views.SoilCalibrationRow(cal, lang, rowStatus).Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (s *Server) renderSoilCalToast(c *gin.Context, lang, status string) {
	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := views.SoilCalibrationToast(lang, status).Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
