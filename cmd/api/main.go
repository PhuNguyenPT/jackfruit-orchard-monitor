package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "GoApp/internal/config"
	"GoApp/internal/server"

	mqtt "github.com/mochi-mqtt/server/v2"
)

var AppVersion = "dev" // overridden by -ldflags in CI/CD

func initLogger(cfg *config.Config) {
	var handler slog.Handler
	opts := &slog.HandlerOptions{Level: cfg.LogLevel}
	if cfg.AppEnv == config.EnvProduction {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	slog.SetDefault(slog.New(handler))
}

func loadMTLSConfig(cfg *config.Config) (*tls.Config, error) {
	caCert, err := os.ReadFile(cfg.TLSCAPath)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)
	cert, err := tls.LoadX509KeyPair(cfg.TLSCertPath, cfg.TLSKeyPath)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		MinVersion:   tls.VersionTLS12,
	}, nil
}

func main() {
	cfg, err := config.LoadAppConfig()
	if err != nil {
		log.Fatalf("invalid config: %v", err) // slog not up yet
	}
	initLogger(cfg)
	cfg.AppVersion = AppVersion
	slog.Info("starting",
		"app", cfg.AppName,
		"version", cfg.AppVersion,
		"env", cfg.AppEnv,
		"log_level", cfg.LogLevel.Level(),
	)

	httpSrv, mqttSrv, err := server.NewServer(cfg)
	if err != nil {
		slog.Error("failed to start server", "err", err)
		os.Exit(1)
	}

	done := make(chan bool, 1)
	go gracefulShutdown(httpSrv, mqttSrv, done)

	if cfg.TLSCertPath != "" {
		tlsCfg, tlsErr := loadMTLSConfig(cfg)
		if tlsErr != nil {
			slog.Error("failed to load mTLS config", "err", tlsErr)
			os.Exit(1)
		}
		httpSrv.TLSConfig = tlsCfg
		slog.Info("mTLS configured, starting HTTPS", "port", cfg.TLSPort)
		err = httpSrv.ListenAndServeTLS(cfg.TLSCertPath, cfg.TLSKeyPath)
	} else {
		slog.Info("starting HTTP", "port", cfg.Port)
		err = httpSrv.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		slog.Error("http server error", "err", err)
		os.Exit(1)
	}
	<-done
	slog.Info("graceful shutdown complete")
}

func gracefulShutdown(httpSrv *http.Server, mqttSrv *mqtt.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	slog.Info("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		slog.Error("HTTP server forced to shutdown", "err", err)
	}
	if err := mqttSrv.Close(); err != nil {
		slog.Error("MQTT broker forced to shutdown", "err", err)
	}

	slog.Info("server exiting")
	done <- true
}
