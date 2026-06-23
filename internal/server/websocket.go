package server

import (
	"net/http"
	"net/url"
	"time"

	config "GoApp/internal/config"
)

const wsHandshakeTimeout = 10 * time.Second

// wsCheckOrigin validates that an incoming WebSocket upgrade request's
// Origin header matches one of the app's configured domains. In dev/test
// environments, localhost and loopback origins are also allowed so the
// dashboard works against a locally-served frontend.
func (s *Server) wsCheckOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return false
	}
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	if s.cfg.AppEnv != config.EnvProduction {
		switch u.Hostname() {
		case "localhost", "127.0.0.1", "::1":
			return true
		}
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
