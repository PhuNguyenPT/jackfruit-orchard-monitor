package server

import (
	"net/http"

	appConfig "GoApp/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes(cfg *appConfig.Config) http.Handler {
	gin.SetMode(cfg.GinMode)
	r := gin.New()
	r.Use(s.nonceMiddleware())
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/api/health", "/.well-known/appspecific/com.chrome.devtools.json"},
	}))
	r.Use(gin.Recovery())
	r.Use(s.resolveUserMiddleware())
	r.Use(s.langMiddleware())

	apiGroup := r.Group("/api")
	apiGroup.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	{
		apiGroup.GET("/", s.apiInfoHandler)
		apiGroup.GET("/health", s.healthHandler)
	}

	r.Static("/public", "./frontend/public")
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./frontend/public/favicon.ico")
	})
	r.GET("/", s.homePageHandler)
	r.HEAD("/", s.homePageHandler)
	r.GET("/contact", s.contactPageHandler)
	r.HEAD("/contact", s.contactPageHandler)
	r.POST("/contact", s.contactFormHandler)
	r.GET("/sitemap.xml", s.sitemapHandler)
	r.HEAD("/sitemap.xml", s.sitemapHandler)
	r.GET("/robots.txt", s.robotsHandler)
	r.HEAD("/robots.txt", s.robotsHandler)
	r.GET("/site.webmanifest", s.webmanifestHandler)
	r.HEAD("/site.webmanifest", s.webmanifestHandler)
	r.GET("/register", s.registerPageHandler)
	r.POST("/register", s.registerHandler)
	r.GET("/login", s.loginPageHandler)
	r.POST("/login", s.loginHandler)
	r.GET("/logout", s.logoutHandler)

	protected := r.Group("/")
	protected.Use(s.authMiddleware())
	{
		protected.GET("/dashboard", s.dashboardPageHandler)
		protected.PUT("/dashboard/name", s.updateUserNameHandler)
		protected.PUT("/dashboard/password", s.updateUserPasswordHandler)
		protected.DELETE("/dashboard/session/:id", s.revokeSessionHandler)
		protected.GET("/sensors", s.sensorsPageHandler)
		protected.GET("/sensors/readings", s.sensorsGridHandler)
		protected.GET("/sensors/ws", s.sensorsWSHandler)
		protected.GET("/sensors/sht40/:addr/history", s.sht40HistoryHandler)
		protected.GET("/sensors/sht40/:addr/ws", s.sht40HistoryWSHandler)
		protected.GET("/sensors/soil/:idx/history", s.soilHistoryHandler)
		protected.GET("/sensors/soil/:idx/ws", s.soilHistoryWSHandler)
	}
	return r
}
