package broker

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"

	"GoApp/internal/database"
)

const (
	topicPrefix = "sht40"
	topicSuffix = "data"
)

type Storage interface {
	InsertSensorReading(ctx context.Context, arg database.InsertSensorReadingParams) error
}

// sensorPayload matches the JSON published by the ESP32 firmware.
type sensorPayload struct {
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
}
type Notifier interface {
	Broadcast(addr string, temperature, humidity float32)
}

type sensorHook struct {
	mqtt.HookBase
	db       Storage
	notifier Notifier
}

func (h *sensorHook) ID() string { return "sensor-hook" }

func (h *sensorHook) Provides(b byte) bool {
	return bytes.Contains([]byte{mqtt.OnPublish}, []byte{b})
}

func (h *sensorHook) OnPublish(_ *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	parts := strings.SplitN(pk.TopicName, "/", 3)
	if len(parts) != 3 || parts[0] != topicPrefix || parts[2] != topicSuffix {
		return pk, nil
	}
	addr := parts[1]

	var sp sensorPayload
	if err := json.Unmarshal(pk.Payload, &sp); err != nil {
		log.Printf("[MQTT] %s/%s/%s: invalid JSON payload: %v — dropping", topicPrefix, addr, topicSuffix, err)
		return pk, nil
	}

	if err := h.db.InsertSensorReading(context.Background(), database.InsertSensorReadingParams{
		Addr:        addr,
		Temperature: sp.Temperature,
		Humidity:    sp.Humidity,
	}); err != nil {
		log.Printf("[MQTT] %s/%s/%s: DB insert failed: %v", topicPrefix, addr, topicSuffix, err)
		return pk, nil
	}

	if h.notifier != nil {
		h.notifier.Broadcast(addr, sp.Temperature, sp.Humidity)
	}
	log.Printf("[MQTT] %s/%s/%s  %.1f %%RH  %.1f °C  → saved", topicPrefix, addr, topicSuffix, sp.Humidity, sp.Temperature)
	return pk, nil
}

// Start launches the embedded MQTT broker on the given TCP port.
func Start(port int, db Storage, notifier Notifier, tlsCfg *tls.Config, user, pass string) (*mqtt.Server, error) {
	server := mqtt.New(&mqtt.Options{})

	if err := server.AddHook(new(auth.Hook), &auth.Options{
		Ledger: &auth.Ledger{
			Auth: auth.AuthRules{
				{Username: auth.RString(user), Password: auth.RString(pass), Allow: true},
			},
		},
	}); err != nil {
		return nil, fmt.Errorf("broker: auth hook: %w", err)
	}
	if err := server.AddHook(&sensorHook{db: db, notifier: notifier}, nil); err != nil {
		return nil, fmt.Errorf("broker: sensor hook: %w", err)
	}

	tcp := listeners.NewTCP(listeners.Config{
		ID:        "tcp8883",
		Address:   fmt.Sprintf(":%d", port),
		TLSConfig: tlsCfg, // nil = plain TCP (dev only)
	})
	if err := server.AddListener(tcp); err != nil {
		return nil, fmt.Errorf("broker: TCP listener: %w", err)
	}

	go func() {
		if err := server.Serve(); err != nil {
			log.Printf("[MQTT] broker stopped: %v", err)
		}
	}()

	log.Printf("[MQTT] broker listening on :%d (TLS=%v)", port, tlsCfg != nil)
	return server, nil
}
