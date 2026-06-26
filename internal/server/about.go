package server

import (
	"log"

	"GoApp/internal/views"

	"github.com/gin-gonic/gin"
)

func (s *Server) aboutPageHandler(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := views.AboutPage(getUserName(c), getLangStr(c), s.siteConfig(c)).Render(c.Request.Context(), c.Writer); err != nil {
		log.Printf("error rendering about page: %v", err)
	}
}
