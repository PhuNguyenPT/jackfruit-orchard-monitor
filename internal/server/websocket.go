package server

import (
	"net/http"
	"net/url"
	"time"
)

const wsHandshakeTimeout = 10 * time.Second

// wsCheckOrigin validates that an incoming WebSocket upgrade request's
// Origin header matches one of the app's configured domains. Used to
// build the shared Upgrader in NewServer.
func (s *Server) wsCheckOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return false
	}
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	for _, base := range s.cfg.BaseURLs {
		bu, err := url.Parse(base)
		if err != nil {
			continue
		}
		if u.Hostname() == bu.Hostname() {
			return true
		}
	}
	return false
}
