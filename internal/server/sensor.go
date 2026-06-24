package server

import (
	"GoApp/internal/database"
	"GoApp/internal/views"
	"log"
	"net/http"
	"strconv"

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
