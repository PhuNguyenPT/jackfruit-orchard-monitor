package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) sitemapHandler(c *gin.Context) {
	base := s.cfg.BaseURL
	body := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    <url>
        <loc>%s/</loc>
        <changefreq>weekly</changefreq>
        <priority>1.0</priority>
    </url>
    <url>
        <loc>%s/contact</loc>
        <changefreq>monthly</changefreq>
        <priority>0.8</priority>
    </url>
</urlset>`, base, base)

	c.Data(http.StatusOK, "application/xml; charset=utf-8", []byte(body))
}

func (s *Server) robotsHandler(c *gin.Context) {
	body := fmt.Sprintf("User-agent: *\nAllow: /\n\nSitemap: %s/sitemap.xml", s.cfg.BaseURL)
	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(body))
}

func (s *Server) webmanifestHandler(c *gin.Context) {
	body := fmt.Sprintf(`{
    "name": "%s",
    "short_name": "%s",
    "icons": [
        {
            "src": "/public/android-chrome-192x192.png",
            "sizes": "192x192",
            "type": "image/png"
        },
        {
            "src": "/public/android-chrome-512x512.png",
            "sizes": "512x512",
            "type": "image/png"
        }
    ],
    "theme_color": "#2563eb",
    "background_color": "#ffffff",
    "display": "standalone"
}`, s.cfg.AppName, s.cfg.AppName)

	c.Data(http.StatusOK, "application/manifest+json; charset=utf-8", []byte(body))
}
