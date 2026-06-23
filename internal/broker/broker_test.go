package broker

import (
	"context"
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"golang.org/x/crypto/bcrypt"

	"GoApp/internal/database"
)

// ---------------------------------------------------------------------------
// Test doubles
// ---------------------------------------------------------------------------

type mockStorage struct {
	// SHT40 Tracking
	airTempHumidCalls []database.InsertAirTempHumidReadingParams
	errOn             int
	callN             int

	// Soil Tracking
	soilCalls []database.InsertSoilMoistureReadingParams
	soilErrOn int
	soilCallN int

	// Auth/ACL Tracking
	getCredErr error
	cred       database.MqttCredential
	acls       []database.MqttAcl
}

func (m *mockStorage) InsertAirTempHumidReading(_ context.Context, arg database.InsertAirTempHumidReadingParams) error {
	m.callN++
	m.airTempHumidCalls = append(m.airTempHumidCalls, arg)
	if m.errOn != 0 && m.callN == m.errOn {
		return errors.New("mock db error")
	}
	return nil
}

func (m *mockStorage) InsertSoilMoistureReading(_ context.Context, arg database.InsertSoilMoistureReadingParams) error {
	m.soilCallN++
	m.soilCalls = append(m.soilCalls, arg)
	if m.soilErrOn != 0 && m.soilCallN == m.soilErrOn {
		return errors.New("mock db error")
	}
	return nil
}
func (m *mockStorage) GetMQTTCredentialByUsername(_ context.Context, username string) (database.MqttCredential, error) {
	if m.getCredErr != nil {
		return database.MqttCredential{}, m.getCredErr
	}
	// Default to returning a matching username to bypass seeding logic smoothly if invoked
	if m.cred.Username == "" {
		return database.MqttCredential{Username: username}, nil
	}
	return m.cred, nil
}

func (m *mockStorage) CreateMQTTCredential(_ context.Context, arg database.CreateMQTTCredentialParams) (database.MqttCredential, error) {
	return database.MqttCredential{Username: arg.Username}, nil
}

func (m *mockStorage) GetMQTTACLByCredentialID(_ context.Context, _ uuid.UUID) ([]database.MqttAcl, error) {
	return m.acls, nil
}

func (m *mockStorage) CreateMQTTACL(_ context.Context, arg database.CreateMQTTACLParams) (database.MqttAcl, error) {
	return database.MqttAcl{CredentialID: arg.CredentialID, Topic: arg.Topic}, nil
}

type airTempHumidCall struct {
	addr        string
	temperature float32
	humidity    float32
	createdAt   time.Time
}

type soilCall struct {
	addr      string
	raw       int
	createdAt time.Time
}

type deviceStatusCall struct {
	clientID  string
	connected bool
}

type mockNotifier struct {
	airTempHumidCalls []airTempHumidCall
	soilCalls         []soilCall
	deviceStatusCalls []deviceStatusCall
}

func (m *mockNotifier) BroadcastAirTempHumid(addr string, temperature, humidity float32, createdAt time.Time) {
	m.airTempHumidCalls = append(m.airTempHumidCalls, airTempHumidCall{addr, temperature, humidity, createdAt})
}

func (m *mockNotifier) BroadcastSoilMoisture(addr string, raw int, createdAt time.Time) {
	m.soilCalls = append(m.soilCalls, soilCall{addr, raw, createdAt})
}

func (m *mockNotifier) BroadcastDeviceStatus(clientID string, connected bool) { // ← add
	m.deviceStatusCalls = append(m.deviceStatusCalls, deviceStatusCall{clientID, connected})
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func makePayload(hum, temp float32) []byte {
	return []byte(fmt.Sprintf(`{"humidity":%.1f,"temperature":%.1f}`, hum, temp))
}

func newTestHook(db Storage, n Notifier) *sensorHook {
	return &sensorHook{db: db, notifier: n}
}

func closeF32(a, b, epsilon float32) bool {
	return float32(math.Abs(float64(a-b))) < epsilon
}

// ---------------------------------------------------------------------------
// sensorHook.ID
// ---------------------------------------------------------------------------

func TestSensorHook_ID(t *testing.T) {
	t.Parallel()
	h := newTestHook(nil, nil)
	if got := h.ID(); got != "sensor-hook" {
		t.Errorf("ID() = %q, want %q", got, "sensor-hook")
	}
}

// ---------------------------------------------------------------------------
// sensorHook.Provides
// ---------------------------------------------------------------------------

func TestSensorHook_Provides(t *testing.T) {
	t.Parallel()
	h := newTestHook(nil, nil)

	for _, b := range []byte{mqtt.OnPublish, mqtt.OnSessionEstablished, mqtt.OnDisconnect} {
		if !h.Provides(b) {
			t.Errorf("Provides(%d) = false, want true", b)
		}
	}
	if h.Provides(mqtt.OnACLCheck) {
		t.Errorf("Provides(OnACLCheck) = true, want false")
	}
}

// ---------------------------------------------------------------------------
// sensorHook.OnPublish — Soil Tests
// ---------------------------------------------------------------------------

func TestSensorHook_OnPublish_Soil(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		topic      string
		payload    []byte
		dbErrOn    int
		wantInsert bool
		wantIdx    int16
		wantRaw    int16
	}{
		{
			name:       "valid soil reading",
			topic:      "mke-s13/1/data",
			payload:    []byte(`{"raw":1500}`),
			wantInsert: true,
			wantIdx:    1,
			wantRaw:    1500,
		},
		{
			name:       "valid soil reading — index 0",
			topic:      "mke-s13/0/data",
			payload:    []byte(`{"raw":3000}`),
			wantInsert: true,
			wantIdx:    0,
			wantRaw:    3000,
		},
		{
			name:       "wrong suffix — drop",
			topic:      "mke-s13/1/status",
			payload:    []byte(`{"raw":1500}`),
			wantInsert: false,
		},
		{
			name:       "non-numeric index — drop",
			topic:      "mke-s13/abc/data",
			payload:    []byte(`{"raw":1500}`),
			wantInsert: false,
		},
		{
			name:       "invalid JSON — drop",
			topic:      "mke-s13/2/data",
			payload:    []byte(`not-json`),
			wantInsert: false,
		},
		{
			name:       "malformed JSON — drop",
			topic:      "mke-s13/2/data",
			payload:    []byte(`{"raw": "bad_value"}`),
			wantInsert: false,
		},
		{
			name:       "db insert error — handle gracefully",
			topic:      "mke-s13/3/data",
			payload:    []byte(`{"raw":2000}`),
			dbErrOn:    1,
			wantInsert: true,
			wantIdx:    3,
			wantRaw:    2000,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := &mockStorage{soilErrOn: tc.dbErrOn}
			notifier := &mockNotifier{}
			h := newTestHook(store, notifier)

			pk := packets.Packet{
				TopicName: tc.topic,
				Payload:   tc.payload,
			}

			got, err := h.OnPublish(nil, pk)

			if err != nil {
				t.Fatalf("OnPublish() returned unexpected error: %v", err)
			}
			if got.TopicName != tc.topic {
				t.Errorf("TopicName modified: got %q, want %q", got.TopicName, tc.topic)
			}

			insertCount := len(store.soilCalls)

			if !tc.wantInsert {
				if insertCount != 0 {
					t.Errorf("expected no DB insert, got %d", insertCount)
				}
				return
			}

			if insertCount != 1 {
				t.Fatalf("expected 1 DB insert, got %d", insertCount)
			}

			if tc.dbErrOn != 0 {
				return // Stop evaluating fields if DB insert simulated an error
			}

			// --- DB insert assertions ---
			arg := store.soilCalls[0]

			if arg.SensorIdx != tc.wantIdx {
				t.Errorf("SensorIdx = %d, want %d", arg.SensorIdx, tc.wantIdx)
			}

			if arg.Raw != tc.wantRaw {
				t.Errorf("Raw = %d, want %d", arg.Raw, tc.wantRaw)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Start — smoke tests
// ---------------------------------------------------------------------------

func TestStart_PlainTCP(t *testing.T) {
	t.Parallel()
	store := &mockStorage{}
	srv, err := Start(18883, store, nil, nil, "testuser", "testpassword")
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if srv == nil {
		t.Fatal("Start() returned nil server")
	}
	t.Cleanup(func() { _ = srv.Close() })
}

func TestStart_DifferentPort(t *testing.T) {
	t.Parallel()
	store := &mockStorage{}
	srv, err := Start(18884, store, nil, nil, "testuser", "testpassword")
	if err != nil {
		t.Fatalf("Start() on alt port error = %v", err)
	}
	t.Cleanup(func() { _ = srv.Close() })
}

func TestSensorHook_OnPublish_SHT40(t *testing.T) {
	t.Parallel()
	store := &mockStorage{}
	notifier := &mockNotifier{}
	h := newTestHook(store, notifier)

	// Using makePayload here!
	payload := makePayload(45.5, 23.2)
	pk := packets.Packet{TopicName: "sht40/12/data", Payload: payload}

	_, err := h.OnPublish(nil, pk)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(notifier.airTempHumidCalls) != 1 {
		t.Fatalf("expected 1 notification broadcast, got %d", len(notifier.airTempHumidCalls))
	}

	// Using closeF32 here to compare floating-point values safely!
	gotTemp := notifier.airTempHumidCalls[0].temperature
	if !closeF32(gotTemp, 23.2, 0.01) {
		t.Errorf("got temperature %f, want 23.2", gotTemp)
	}
}

// ---------------------------------------------------------------------------
// mqttTopicMatch
// ---------------------------------------------------------------------------

func TestMqttTopicMatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		pattern string
		topic   string
		want    bool
	}{
		{"exact single-level wildcard", "sht40/+/data", "sht40/1/data", true},
		{"exact single-level wildcard — different addr", "sht40/+/data", "sht40/12/data", true},
		{"mismatched prefix", "sht40/+/data", "mke-s13/1/data", false},
		{"too many levels for + pattern", "sht40/+/data", "sht40/1/extra/data", false},
		{"too few levels for + pattern", "sht40/+/data", "sht40/1", false},
		{"exact match no wildcards", "sht40/1/data", "sht40/1/data", true},
		{"exact match no wildcards — mismatch", "sht40/1/data", "sht40/2/data", false},

		// '#' as the final (only valid) position
		{"# matches single extra level", "sht40/#", "sht40/1", true},
		{"# matches multiple extra levels", "sht40/#", "sht40/1/data", true},
		{"# matches deeply nested", "sht40/#", "sht40/1/data/extra", true},
		{"# alone matches everything", "#", "anything/at/all", true},
		{"# alone matches single level", "#", "anything", true},

		// '#' in a non-final position is malformed per MQTT 3.1.1 §4.7 — must never match
		{"# mid-pattern — malformed, no match", "sht40/#/data", "sht40/1/data", false},
		{"# first of two — malformed, no match", "#/data", "1/data", false},
		{"# first of three — malformed, no match", "#/foo/bar", "1/foo/bar", false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := mqttTopicMatch(tc.pattern, tc.topic); got != tc.want {
				t.Errorf("mqttTopicMatch(%q, %q) = %v, want %v", tc.pattern, tc.topic, got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// authHook.ID / Provides
// ---------------------------------------------------------------------------

func TestAuthHook_ID(t *testing.T) {
	t.Parallel()
	h := &authHook{}
	if got := h.ID(); got != "auth-ledger" {
		t.Errorf("ID() = %q, want %q", got, "auth-ledger")
	}
}

func TestAuthHook_Provides(t *testing.T) {
	t.Parallel()
	h := &authHook{}

	for _, b := range []byte{mqtt.OnConnectAuthenticate, mqtt.OnACLCheck} {
		if !h.Provides(b) {
			t.Errorf("Provides(%d) = false, want true", b)
		}
	}
	for _, b := range []byte{mqtt.OnPublish, mqtt.OnDisconnect, mqtt.OnSessionEstablished} {
		if h.Provides(b) {
			t.Errorf("Provides(%d) = true, want false", b)
		}
	}
}

// ---------------------------------------------------------------------------
// authHook.OnConnectAuthenticate
// ---------------------------------------------------------------------------

func TestAuthHook_OnConnectAuthenticate(t *testing.T) {
	t.Parallel()

	t.Run("valid credentials authenticate", func(t *testing.T) {
		t.Parallel()
		hash, err := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
		if err != nil {
			t.Fatalf("failed to hash password: %v", err)
		}
		store := &mockStorage{cred: database.MqttCredential{Username: "esp32", Password: string(hash)}}
		h := &authHook{db: store}

		pk := packets.Packet{Connect: packets.ConnectParams{
			Username: []byte("esp32"),
			Password: []byte("correctpass"),
		}}

		if !h.OnConnectAuthenticate(nil, pk) {
			t.Error("OnConnectAuthenticate() = false, want true for correct credentials")
		}
	})

	t.Run("wrong password rejected", func(t *testing.T) {
		t.Parallel()
		hash, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
		store := &mockStorage{cred: database.MqttCredential{Username: "esp32", Password: string(hash)}}
		h := &authHook{db: store}

		pk := packets.Packet{Connect: packets.ConnectParams{
			Username: []byte("esp32"),
			Password: []byte("wrongpass"),
		}}

		if h.OnConnectAuthenticate(nil, pk) {
			t.Error("OnConnectAuthenticate() = true, want false for wrong password")
		}
	})

	t.Run("unknown user rejected", func(t *testing.T) {
		t.Parallel()
		store := &mockStorage{getCredErr: errors.New("not found")}
		h := &authHook{db: store}

		pk := packets.Packet{Connect: packets.ConnectParams{
			Username: []byte("ghost"),
			Password: []byte("whatever"),
		}}

		if h.OnConnectAuthenticate(nil, pk) {
			t.Error("OnConnectAuthenticate() = true, want false for lookup error")
		}
	})
}

// ---------------------------------------------------------------------------
// authHook.OnACLCheck
// ---------------------------------------------------------------------------

func TestAuthHook_OnACLCheck(t *testing.T) {
	t.Parallel()

	t.Run("topic covered by ACL is allowed", func(t *testing.T) {
		t.Parallel()
		store := &mockStorage{
			cred: database.MqttCredential{Username: "esp32"},
			acls: []database.MqttAcl{{Topic: "sht40/+/data"}, {Topic: "mke-s13/+/data"}},
		}
		h := &authHook{db: store}
		cl := &mqtt.Client{Properties: mqtt.ClientProperties{Username: []byte("esp32")}}

		if !h.OnACLCheck(cl, "sht40/1/data", true) {
			t.Error("OnACLCheck() = false, want true for topic matching an ACL")
		}
	})

	t.Run("topic not covered by any ACL is denied", func(t *testing.T) {
		t.Parallel()
		store := &mockStorage{
			cred: database.MqttCredential{Username: "esp32"},
			acls: []database.MqttAcl{{Topic: "sht40/+/data"}},
		}
		h := &authHook{db: store}
		cl := &mqtt.Client{Properties: mqtt.ClientProperties{Username: []byte("esp32")}}

		if h.OnACLCheck(cl, "admin/shutdown", true) {
			t.Error("OnACLCheck() = true, want false for topic not covered by any ACL")
		}
	})

	t.Run("credential lookup failure denies", func(t *testing.T) {
		t.Parallel()
		store := &mockStorage{getCredErr: errors.New("lookup failed")}
		h := &authHook{db: store}
		cl := &mqtt.Client{Properties: mqtt.ClientProperties{Username: []byte("ghost")}}

		if h.OnACLCheck(cl, "sht40/1/data", true) {
			t.Error("OnACLCheck() = true, want false when credential lookup fails")
		}
	})
}

// ---------------------------------------------------------------------------
// sensorHook.OnSessionEstablished / OnDisconnect
// ---------------------------------------------------------------------------

func makeClient(username string) *mqtt.Client {
	cl := &mqtt.Client{}
	cl.Properties.Username = []byte(username)
	return cl
}

func TestSensorHook_OnSessionEstablished_NotifiesConnected(t *testing.T) {
	t.Parallel()
	notifier := &mockNotifier{}
	h := newTestHook(nil, notifier)

	h.OnSessionEstablished(makeClient("esp32"), packets.Packet{})

	if len(notifier.deviceStatusCalls) != 1 {
		t.Fatalf("expected 1 status call, got %d", len(notifier.deviceStatusCalls))
	}
	call := notifier.deviceStatusCalls[0]
	if call.clientID != "esp32" {
		t.Errorf("clientID = %q, want %q", call.clientID, "esp32")
	}
	if !call.connected {
		t.Error("connected = false, want true")
	}
}

func TestSensorHook_OnDisconnect_NotifiesDisconnected(t *testing.T) {
	t.Parallel()
	notifier := &mockNotifier{}
	h := newTestHook(nil, notifier)

	h.OnDisconnect(makeClient("esp32"), nil, false)

	if len(notifier.deviceStatusCalls) != 1 {
		t.Fatalf("expected 1 status call, got %d", len(notifier.deviceStatusCalls))
	}
	call := notifier.deviceStatusCalls[0]
	if call.clientID != "esp32" {
		t.Errorf("clientID = %q, want %q", call.clientID, "esp32")
	}
	if call.connected {
		t.Error("connected = true, want false")
	}
}

func TestSensorHook_OnDisconnect_NilNotifier_NoPanic(t *testing.T) {
	t.Parallel()
	h := newTestHook(nil, nil) // notifier is nil
	// Must not panic
	h.OnDisconnect(makeClient("esp32"), nil, false)
}

func TestSensorHook_OnSessionEstablished_NilNotifier_NoPanic(t *testing.T) {
	t.Parallel()
	h := newTestHook(nil, nil)
	h.OnSessionEstablished(makeClient("esp32"), packets.Packet{})
}
