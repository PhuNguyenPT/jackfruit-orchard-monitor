package server

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

// requestBaseURL returns the scheme+host the current request actually
// arrived on, provided it matches one of the configured BaseURLs.
// Falls back to the primary (first) BaseURL if the request's Host
// doesn't match anything known — this should only happen for malformed
// or unexpected Host headers.
func (s *Server) requestBaseURL(c *gin.Context) string {
	if len(s.cfg.BaseURLs) == 0 {
		return "" // no configured domain to fall back to
	}
	host := c.Request.Host

	for _, base := range s.cfg.BaseURLs {
		bu, err := url.Parse(base)
		if err != nil {
			continue
		}
		if bu.Host == host || bu.Hostname() == host {
			return base
		}
	}

	// Unrecognized Host header — fall back to primary domain.
	return s.cfg.BaseURLs[0]
}
