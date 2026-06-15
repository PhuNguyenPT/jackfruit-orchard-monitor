package server

import (
	"GoApp/internal/ctxutil"
	"context"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session_token")
		if err != nil {
			next := url.QueryEscape(c.Request.URL.RequestURI())
			c.Redirect(http.StatusFound, "/login?next="+next)
			c.Abort()
			return
		}

		session, err := s.db.GetSessionByToken(c.Request.Context(), token)
		if err != nil {
			secure := s.cfg.AppEnv == EnvProduction
			c.SetCookie("session_token", "", -1, "/", "", secure, true)
			next := url.QueryEscape(c.Request.URL.RequestURI())
			c.Redirect(http.StatusFound, "/login?next="+next)
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
		nonce := c.GetHeader("X-Nonce")

		if nonce == "" {
			nonce = "dev-fallback-nonce"
		}

		c.Set("nonce", nonce)

		ctx := context.WithValue(c.Request.Context(), ctxutil.NonceKey, nonce)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
