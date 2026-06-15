package broker

import (
	"context"
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/google/uuid"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"

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

type mockNotifier struct {
	airTempHumidCalls []struct {
		addr        string
		temperature float32
		humidity    float32
	}
	soilCalls []struct {
		addr string
		raw  int
	}
}

func (m *mockNotifier) BroadcastAirTempHumid(addr string, temperature, humidity float32) {
	m.airTempHumidCalls = append(m.airTempHumidCalls, struct {
		addr        string
		temperature float32
		humidity    float32
	}{addr, temperature, humidity})
}

func (m *mockNotifier) BroadcastSoilMoisture(addr string, raw int) {
	m.soilCalls = append(m.soilCalls, struct {
		addr string
		raw  int
	}{addr, raw})
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

	if !h.Provides(mqtt.OnPublish) {
		t.Error("Provides(mqtt.OnPublish) = false, want true")
	}
	for _, b := range []byte{mqtt.OnPublish - 1, mqtt.OnPublish + 1} {
		if h.Provides(b) {
			t.Errorf("Provides(%d) = true, want false", b)
		}
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
