package server

import (
	"GoApp/internal/database"
	"GoApp/internal/views"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func (s *Server) sensorsPageHandler(c *gin.Context) {
	lang := getLangStr(c)

	// 1. Fetch SHT40 Data
	shtReadings, err := s.db.GetLatestAirTempHumidReadings(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 2. Fetch MKE-S13 Data
	soilReadings, err := s.db.GetLatestSoilMoistureReadings(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	deviceStatuses := s.hub.GetDeviceStatuses()

	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := views.SensorsPage(shtReadings, soilReadings, lang, getUserName(c), s.cfg.SoilDryValue, s.cfg.SoilWetValue, s.siteConfig(c), deviceStatuses).Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (s *Server) sensorsGridHandler(c *gin.Context) {
	lang := getLangStr(c)

	shtReadings, err := s.db.GetLatestAirTempHumidReadings(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	soilReadings, err := s.db.GetLatestSoilMoistureReadings(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := views.SensorGrid(shtReadings, soilReadings, lang, s.cfg.SoilDryValue, s.cfg.SoilWetValue).Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (s *Server) sensorsWSHandler(c *gin.Context) {
	lang := getLangStr(c)
	conn, err := s.wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WS] upgrade error: %v", err)
		return
	}

	s.hub.register(conn, lang)
	defer s.hub.unregister(conn)
	s.hub.pushDeviceStatusesToClient(conn, lang)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseNormalClosure,
			) {
				log.Printf("[WS] unexpected close: %v", err)
			}
			break
		}
	}
}

func (s *Server) sht40HistoryHandler(c *gin.Context) {
	lang := getLangStr(c)
	addr, err := strconv.ParseInt(c.Param("addr"), 10, 16)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	rows, err := s.db.GetAirTempHumidReadingsByAddr(c.Request.Context(),
		database.GetAirTempHumidReadingsByAddrParams{
			Addr:  int16(addr),
			Limit: 100,
		})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := views.SHT40HistoryPage(rows, int16(addr), lang, s.siteConfig(c)).
		Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (s *Server) soilHistoryHandler(c *gin.Context) {
	lang := getLangStr(c)
	idx, err := strconv.ParseInt(c.Param("idx"), 10, 16)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	rows, err := s.db.GetSoilMoistureReadingsBySensorIdx(c.Request.Context(),
		database.GetSoilMoistureReadingsBySensorIdxParams{
			SensorIdx: int16(idx),
			Limit:     100,
		})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := views.SoilHistoryPage(rows, int16(idx), lang, s.siteConfig(c), s.cfg.SoilDryValue, s.cfg.SoilWetValue).
		Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func parseChartTimestamp(s, lang string) (time.Time, error) {
	format := "02-01 15:04:05"
	if lang != "vi" {
		format = "01-02 15:04:05"
	}
	parsed, err := time.ParseInLocation(format, s, hubVietnamTZ)
	if err != nil {
		return time.Time{}, err
	}
	now := time.Now().In(hubVietnamTZ)
	result := time.Date(now.Year(), parsed.Month(), parsed.Day(),
		parsed.Hour(), parsed.Minute(), parsed.Second(), 0, hubVietnamTZ)
	if result.After(now) {
		result = result.AddDate(-1, 0, 0)
	}
	return result, nil
}

func (s *Server) sht40HistoryWSHandler(c *gin.Context) {
	lang := getLangStr(c)

	addr, err := strconv.ParseInt(c.Param("addr"), 10, 16)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	conn, err := s.wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WS] sht40 history upgrade error: %v", err)
		return
	}

	s.hub.registerSHT40Chart(conn, int16(addr), lang)
	defer s.hub.unregisterSHT40Chart(conn)

	if sinceStr := c.Query("since"); sinceStr != "" {
		if since, perr := parseChartTimestamp(sinceStr, lang); perr == nil {
			s.hub.pushSHT40Backfill(c.Request.Context(), s.db, conn, int16(addr), since, lang)
		} else {
			log.Printf("[WS] sht40 since parse error: %v", perr)
		}
	}

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("[WS] sht40 history unexpected close: %v", err)
			}
			break
		}
	}
}

func (s *Server) soilHistoryWSHandler(c *gin.Context) {
	lang := getLangStr(c)
	idx, err := strconv.ParseInt(c.Param("idx"), 10, 16)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	conn, err := s.wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WS] soil history upgrade error: %v", err)
		return
	}

	s.hub.registerSoilChart(conn, int16(idx), lang)
	defer s.hub.unregisterSoilChart(conn)

	if sinceStr := c.Query("since"); sinceStr != "" {
		if since, perr := parseChartTimestamp(sinceStr, lang); perr == nil {
			s.hub.pushSoilBackfill(c.Request.Context(), s.db, conn, int16(idx), since, lang)
		} else {
			log.Printf("[WS] soil since parse error: %v", perr)
		}
	}

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("[WS] soil history unexpected close: %v", err)
			}
			break
		}
	}
}
