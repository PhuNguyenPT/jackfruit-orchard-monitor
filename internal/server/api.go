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
	c.JSON(http.StatusOK, s.db.Health())
}
