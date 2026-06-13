package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"GoApp/internal/server"

	mqtt "github.com/mochi-mqtt/server/v2"
)

func loadMTLSConfig(cfg *server.Config) (*tls.Config, error) {
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

func gracefulShutdown(httpSrv *http.Server, mqttSrv *mqtt.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server forced to shutdown: %v", err)
	}
	if err := mqttSrv.Close(); err != nil {
		log.Printf("MQTT broker forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
	done <- true
}

func main() {
	cfg, err := server.LoadConfig()
	if err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	httpSrv, mqttSrv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	done := make(chan bool, 1)
	go gracefulShutdown(httpSrv, mqttSrv, done)

	if cfg.TLSCertPath != "" {
		tlsCfg, tlsErr := loadMTLSConfig(cfg)
		if tlsErr != nil {
			log.Fatalf("failed to load mTLS config: %v", tlsErr)
		}
		httpSrv.TLSConfig = tlsCfg
		log.Printf("mTLS configured, starting HTTPS on %d", cfg.TLSPort)
		err = httpSrv.ListenAndServeTLS(cfg.TLSCertPath, cfg.TLSKeyPath)
	} else {
		log.Printf("starting HTTP on %d", cfg.Port)
		err = httpSrv.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}
	<-done
	log.Println("Graceful shutdown complete.")
}
