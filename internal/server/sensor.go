package server

import (
	"log"
	"net/http"

	"GoApp/internal/views"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *Server) sensorsPageHandler(c *gin.Context) {
	lang := getLangStr(c)
	readings, err := s.db.GetLatestReadings(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := views.SensorsPage(readings, lang, getUserName(c)).Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (s *Server) sensorsGridHandler(c *gin.Context) {
	lang := getLangStr(c)
	readings, err := s.db.GetLatestReadings(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := views.SensorGrid(readings, lang).Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (s *Server) sensorsWSHandler(c *gin.Context) {
	lang := getLangStr(c)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WS] upgrade error: %v", err)
		return
	}

	s.hub.register(conn, lang)
	defer s.hub.unregister(conn)

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
