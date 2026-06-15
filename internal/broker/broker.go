package broker

import (
	"bytes"
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/google/uuid"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
	"golang.org/x/crypto/bcrypt"

	"GoApp/internal/database"
)

const (
	topicWildcard = "+"
	topicSuffix   = "data"
)

// Centralized list of all supported sensor prefixes
var validPrefixes = []string{"sht40", "mke-s13"}

type Storage interface {
	InsertAirTempHumidReading(ctx context.Context, arg database.InsertAirTempHumidReadingParams) error
	InsertSoilMoistureReading(ctx context.Context, arg database.InsertSoilMoistureReadingParams) error
	GetMQTTCredentialByUsername(ctx context.Context, username string) (database.MqttCredential, error)
	CreateMQTTCredential(ctx context.Context, arg database.CreateMQTTCredentialParams) (database.MqttCredential, error)
	GetMQTTACLByCredentialID(ctx context.Context, credentialID uuid.UUID) ([]database.MqttAcl, error)
	CreateMQTTACL(ctx context.Context, arg database.CreateMQTTACLParams) (database.MqttAcl, error)
}

// sensorPayload matches the JSON published by the ESP32 firmware.
type sensorPayload struct {
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
}
type Notifier interface {
	BroadcastAirTempHumid(addr string, temperature, humidity float32)
	BroadcastSoilMoisture(addr string, raw int)
}

type sensorHook struct {
	mqtt.HookBase
	db       Storage
	notifier Notifier
}

func seedCredential(db Storage, user, pass string) error {
	ctx := context.Background()
	cred, err := db.GetMQTTCredentialByUsername(ctx, user)

	if errors.Is(err, sql.ErrNoRows) {
		// User doesn't exist yet, create a fresh one
		hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("seedCredential: hash: %w", err)
		}
		cred, err = db.CreateMQTTCredential(ctx, database.CreateMQTTCredentialParams{
			Username: user,
			Password: string(hash),
		})
		if err != nil {
			return fmt.Errorf("seedCredential: insert: %w", err)
		}
		log.Printf("[MQTT] created fresh credential for %q", user)
	} else if err != nil {
		return fmt.Errorf("seedCredential: lookup: %w", err)
	} else {
		log.Printf("[MQTT] credential for %q already in DB, verifying ACLs...", user)
	}

	// Always fetch existing ACLs to see what's missing
	existingACLs, err := db.GetMQTTACLByCredentialID(ctx, cred.ID)
	if err != nil {
		return fmt.Errorf("seedCredential: get acls: %w", err)
	}

	aclMap := make(map[string]bool)
	for _, acl := range existingACLs {
		aclMap[acl.Topic] = true
	}

	seededCount := 0
	for _, prefix := range validPrefixes {
		topic := fmt.Sprintf("%s/%s/%s", prefix, topicWildcard, topicSuffix)

		// Skip if this specific ACL row already survived
		if aclMap[topic] {
			continue
		}

		if _, err := db.CreateMQTTACL(ctx, database.CreateMQTTACLParams{
			CredentialID: cred.ID,
			Topic:        topic,
			Permission:   "rw",
		}); err != nil {
			return fmt.Errorf("seedCredential: acl %q: %w", topic, err)
		}
		seededCount++
	}

	log.Printf("[MQTT] seeding complete. Added %d missing ACL topics.", seededCount)
	return nil
}

type authHook struct {
	mqtt.HookBase
	db Storage
}

func (h *authHook) ID() string { return "auth-ledger" }

func (h *authHook) Provides(b byte) bool {
	return bytes.Contains([]byte{mqtt.OnConnectAuthenticate, mqtt.OnACLCheck}, []byte{b})
}

func (h *authHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	cred, err := h.db.GetMQTTCredentialByUsername(context.Background(), string(pk.Connect.Username))
	if err != nil {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(cred.Password), []byte(pk.Connect.Password)) == nil
}

func (h *authHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	cred, err := h.db.GetMQTTCredentialByUsername(context.Background(), string(cl.Properties.Username))
	if err != nil {
		return false
	}
	acls, err := h.db.GetMQTTACLByCredentialID(context.Background(), cred.ID)
	if err != nil {
		return false
	}
	for _, acl := range acls {
		if mqttTopicMatch(acl.Topic, topic) {
			return true
		}
	}
	return false
}

// mqttTopicMatch supports + (single level) and # (multi level) wildcards.
func mqttTopicMatch(pattern, topic string) bool {
	pp := strings.Split(pattern, "/")
	tp := strings.Split(topic, "/")
	for i, p := range pp {
		if p == "#" {
			return true
		}
		if i >= len(tp) {
			return false
		}
		if p != "+" && p != tp[i] {
			return false
		}
	}
	return len(pp) == len(tp)
}

func (h *sensorHook) ID() string { return "sensor-hook" }

func (h *sensorHook) Provides(b byte) bool {
	return bytes.Contains([]byte{mqtt.OnPublish}, []byte{b})
}

func (h *sensorHook) OnPublish(_ *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	parts := strings.SplitN(pk.TopicName, "/", 3)
	// Only validate that we have 3 parts and it ends with the correct suffix
	if len(parts) != 3 || parts[2] != topicSuffix {
		return pk, nil
	}

	prefix := parts[0]
	addr := parts[1]

	addrInt, err := strconv.ParseInt(addr, 10, 16)
	if err != nil {
		log.Printf("[MQTT] %s: invalid sensor address: %v — dropping", pk.TopicName, err)
		return pk, nil
	}

	switch prefix {
	case "sht40":
		var sp sensorPayload
		if err := json.Unmarshal(pk.Payload, &sp); err != nil {
			log.Printf("[MQTT] %s: invalid JSON payload: %v — dropping", pk.TopicName, err)
			return pk, nil
		}

		if err := h.db.InsertAirTempHumidReading(context.Background(), database.InsertAirTempHumidReadingParams{
			Addr:        int16(addrInt),
			Temperature: int16(math.Round(float64(sp.Temperature) * 10)),
			Humidity:    int16(math.Round(float64(sp.Humidity) * 10)),
		}); err != nil {
			log.Printf("[MQTT] %s: DB insert failed: %v", pk.TopicName, err)
			return pk, nil
		}

		if h.notifier != nil {
			h.notifier.BroadcastAirTempHumid(addr, sp.Temperature, sp.Humidity)
		}
		log.Printf("[MQTT] %s  %.1f %%RH  %.1f °C  → saved", pk.TopicName, sp.Humidity, sp.Temperature)

	case "mke-s13":
		var sp struct {
			Raw int16 `json:"raw"`
		}
		if err := json.Unmarshal(pk.Payload, &sp); err != nil {
			log.Printf("[MQTT] %s: invalid JSON payload: %v — dropping", pk.TopicName, err)
			return pk, nil
		}

		if err := h.db.InsertSoilMoistureReading(context.Background(), database.InsertSoilMoistureReadingParams{
			SensorIdx: int16(addrInt),
			Raw:       sp.Raw,
		}); err != nil {
			log.Printf("[MQTT] %s: DB insert failed: %v", pk.TopicName, err)
			return pk, nil
		}

		if h.notifier != nil {
			h.notifier.BroadcastSoilMoisture(addr, int(sp.Raw))
		}

		log.Printf("[MQTT] %s  %d raw  → saved", pk.TopicName, sp.Raw)

	default:
		// Silently drop unrecognized prefixes
		return pk, nil
	}

	return pk, nil
}

// Start launches the embedded MQTT broker on the given TCP port.
func Start(port int, db Storage, notifier Notifier, tlsCfg *tls.Config, user, pass string) (*mqtt.Server, error) {
	if err := seedCredential(db, user, pass); err != nil {
		return nil, fmt.Errorf("broker: seed credentials: %w", err)
	}

	server := mqtt.New(&mqtt.Options{})

	if err := server.AddHook(&authHook{db: db}, nil); err != nil {
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
