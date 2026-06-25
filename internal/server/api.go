package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) apiInfoHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Go Server API",
		"version": s.cfg.AppVersion,
		"endpoints": gin.H{
			"health": "/api/health",
		},
	})
}

func (s *Server) healthHandler(c *gin.Context) {
	health := s.db.Health()

	resp := gin.H{
		"status":  health["status"],
		"message": health["message"],
		"version": s.cfg.AppVersion,
		"env":     s.cfg.AppEnv,
	}

	switch health["status"] {
	case "up":
		c.JSON(http.StatusOK, resp)
	case "degraded":
		c.JSON(http.StatusOK, resp) // alive but warn
	case "down":
		c.JSON(http.StatusServiceUnavailable, resp)
	default:
		c.JSON(http.StatusServiceUnavailable, resp)
	}
}
