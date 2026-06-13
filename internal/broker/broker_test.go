package broker

import (
	"context"
	"errors"
	"fmt"
	"math"
	"testing"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"

	"GoApp/internal/database"
)

// ---------------------------------------------------------------------------
// Test doubles
// ---------------------------------------------------------------------------

type mockStorage struct {
	calls []database.InsertSensorReadingParams
	errOn int
	callN int
}

func (m *mockStorage) InsertSensorReading(_ context.Context, arg database.InsertSensorReadingParams) error {
	m.callN++
	m.calls = append(m.calls, arg)
	if m.errOn != 0 && m.callN == m.errOn {
		return errors.New("mock db error")
	}
	return nil
}

type mockNotifier struct {
	calls []struct {
		addr        string
		temperature float32
		humidity    float32
	}
}

func (m *mockNotifier) Broadcast(addr string, temperature, humidity float32) {
	m.calls = append(m.calls, struct {
		addr        string
		temperature float32
		humidity    float32
	}{addr, temperature, humidity})
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
// sensorHook.OnPublish — table-driven
// ---------------------------------------------------------------------------

func TestSensorHook_OnPublish(t *testing.T) {
	t.Parallel()

	const eps = float32(0.05)

	tests := []struct {
		name          string
		topic         string
		payload       []byte
		dbErrOn       int
		wantInsert    bool
		wantBroadcast bool
		wantAddr      string
		wantHum       float32
		wantTemp      float32
	}{
		{
			name:          "valid reading — positive temperature",
			topic:         "sht40/sensor1/data",
			payload:       makePayload(55.3, 27.4),
			wantInsert:    true,
			wantBroadcast: true,
			wantAddr:      "sensor1",
			wantHum:       55.3,
			wantTemp:      27.4,
		},
		{
			name:          "valid reading — negative temperature",
			topic:         "sht40/sensor3/data",
			payload:       makePayload(80.0, -5.2),
			wantInsert:    true,
			wantBroadcast: true,
			wantAddr:      "sensor3",
			wantHum:       80.0,
			wantTemp:      -5.2,
		},
		{
			name:          "valid reading — zero temperature boundary",
			topic:         "sht40/sensor10/data",
			payload:       makePayload(40.0, 0.0),
			wantInsert:    true,
			wantBroadcast: true,
			wantAddr:      "sensor10",
			wantHum:       40.0,
			wantTemp:      0.0,
		},
		{
			name:       "non-matching topic — pass through",
			topic:      "other/topic/foo",
			payload:    makePayload(50.0, 25.0),
			wantInsert: false,
		},
		{
			name:       "bare prefix only — pass through",
			topic:      "sht40",
			payload:    makePayload(50.0, 25.0),
			wantInsert: false,
		},
		{
			name:       "wrong prefix — pass through",
			topic:      "device/sht40/sensor1/data",
			payload:    makePayload(50.0, 25.0),
			wantInsert: false,
		},
		{
			name:       "wrong suffix — pass through",
			topic:      "sht40/sensor1/status",
			payload:    makePayload(50.0, 25.0),
			wantInsert: false,
		},
		{
			name:       "invalid JSON — drop",
			topic:      "sht40/sensor2/data",
			payload:    []byte("not-json"),
			wantInsert: false,
		},
		{
			name:       "empty payload — drop",
			topic:      "sht40/sensor2/data",
			payload:    []byte{},
			wantInsert: false,
		},
		{
			name:       "malformed JSON — drop",
			topic:      "sht40/sensor2/data",
			payload:    []byte(`{"temperature": "bad_value"}`),
			wantInsert: false,
		},
		{
			name:          "db insert error — no broadcast",
			topic:         "sht40/sensor5/data",
			payload:       makePayload(60.0, 30.0),
			dbErrOn:       1,
			wantInsert:    true,
			wantBroadcast: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := &mockStorage{errOn: tc.dbErrOn}
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

			insertCount := len(store.calls)
			broadcastCount := len(notifier.calls)

			if !tc.wantInsert {
				if insertCount != 0 {
					t.Errorf("expected no DB insert, got %d", insertCount)
				}
				if broadcastCount != 0 {
					t.Errorf("expected no broadcast, got %d", broadcastCount)
				}
				return
			}

			if insertCount != 1 {
				t.Fatalf("expected 1 DB insert, got %d", insertCount)
			}

			if tc.wantBroadcast && broadcastCount != 1 {
				t.Errorf("expected 1 broadcast, got %d", broadcastCount)
			}
			if !tc.wantBroadcast && broadcastCount != 0 {
				t.Errorf("expected no broadcast, got %d", broadcastCount)
			}

			if tc.dbErrOn != 0 {
				return
			}

			arg := store.calls[0]
			if arg.Addr != tc.wantAddr {
				t.Errorf("Addr = %q, want %q", arg.Addr, tc.wantAddr)
			}
			if !closeF32(arg.Humidity, tc.wantHum, eps) {
				t.Errorf("Humidity = %.2f, want %.2f", arg.Humidity, tc.wantHum)
			}
			if !closeF32(arg.Temperature, tc.wantTemp, eps) {
				t.Errorf("Temperature = %.2f, want %.2f", arg.Temperature, tc.wantTemp)
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
