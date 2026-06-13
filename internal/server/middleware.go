package server

import (
	"GoApp/internal/ctxutil"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session_token")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		session, err := s.db.GetSessionByToken(c.Request.Context(), token)
		if err != nil {
			secure := s.cfg.AppEnv == EnvProduction
			c.SetCookie("session_token", "", -1, "/", "", secure, true)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Set("userID", session.UserID)
		c.Next()
	}
}

func (s *Server) resolveUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session_token")
		if err != nil {
			c.Set("userName", "")
			c.Next()
			return
		}

		session, err := s.db.GetSessionByToken(c.Request.Context(), token)
		if err != nil {
			c.Set("userName", "")
			c.Next()
			return
		}

		user, err := s.db.GetUserByID(c.Request.Context(), session.UserID)
		if err != nil {
			log.Printf("resolveUserMiddleware: user not found: %v", err)
			c.Set("userName", "")
			c.Next()
			return
		}

		c.Set("userName", user.Name)
		c.Next()
	}
}

func getUserName(c *gin.Context) string {
	name, _ := c.Get("userName")
	if s, ok := name.(string); ok {
		return s
	}
	return ""
}

func (s *Server) nonceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get the nonce passed from Nginx
		nonce := c.GetHeader("X-Nonce")

		// 2. Fallback for local development when bypassing Nginx
		if nonce == "" {
			nonce = "dev-fallback-nonce"
		}

		// 3. Optional: Store in Gin's context if you need it in Gin handlers
		c.Set("nonce", nonce)

		// 4. CRITICAL: Store in the underlying standard request Context for Templ!
		ctx := context.WithValue(c.Request.Context(), ctxutil.NonceKey, nonce)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
