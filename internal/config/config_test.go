package config

import (
	"log/slog"
	"maps"
	"os"
	"testing"
)

func setEnv(t *testing.T, env map[string]string) {
	t.Helper()
	for k, v := range env {
		os.Setenv(k, v)
	}
	t.Cleanup(func() {
		for k := range env {
			os.Unsetenv(k)
		}
	})
}

var validEnv = map[string]string{
	"PORT":              "8080",
	"APP_ENV":           "test",
	"GIN_MODE":          "debug",
	"POSTGRES_HOST":     "localhost",
	"POSTGRES_PORT":     "5432",
	"POSTGRES_DATABASE": "mydb",
	"POSTGRES_USERNAME": "myuser",
	"POSTGRES_PASSWORD": "mypassword",
	"POSTGRES_SCHEMA":   "public",
	"MQTT_PASS":         "testmqttpass",
}

func TestLoadAppConfig_Valid(t *testing.T) {
	setEnv(t, validEnv)

	cfg, err := LoadAppConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %v", cfg.Port)
	}
	if cfg.GinMode != "debug" {
		t.Errorf("expected gin mode debug, got %v", cfg.GinMode)
	}
}

func TestLoadAppConfig_MissingRequired(t *testing.T) {
	setEnv(t, validEnv)
	os.Unsetenv("APP_ENV") // no default for this one

	_, err := LoadAppConfig()
	if err == nil {
		t.Error("expected error for missing APP_ENV, got nil")
	}
}
func TestLoadAppConfig_InvalidGinMode(t *testing.T) {
	env := maps.Clone(validEnv)
	env["GIN_MODE"] = "invalid"
	setEnv(t, env)

	_, err := LoadAppConfig()
	if err == nil {
		t.Error("expected error for invalid GIN_MODE, got nil")
	}
}

func TestLoadAppConfig_ShortPassword(t *testing.T) {
	env := maps.Clone(validEnv)
	env["POSTGRES_PASSWORD"] = "short"
	setEnv(t, env)

	_, err := LoadAppConfig()
	if err == nil {
		t.Error("expected error for password shorter than 8 chars, got nil")
	}
}

func TestLoadAppConfig_NonNumericDBPort(t *testing.T) {
	env := maps.Clone(validEnv)
	env["POSTGRES_PORT"] = "notaport"
	setEnv(t, env)

	_, err := LoadAppConfig()
	if err == nil {
		t.Error("expected error for non-numeric POSTGRES_PORT, got nil")
	}
}

func TestLoadAppConfig_ShortMQTTPassword(t *testing.T) {
	env := maps.Clone(validEnv)
	env["MQTT_PASS"] = "short"
	setEnv(t, env)

	_, err := LoadAppConfig()
	if err == nil {
		t.Error("expected error for MQTT_PASS shorter than 8 chars, got nil")
	}
}

func TestLoadAppConfig_MissingMQTTPass(t *testing.T) {
	env := maps.Clone(validEnv)
	delete(env, "MQTT_PASS")
	setEnv(t, env)

	_, err := LoadAppConfig()
	if err == nil {
		t.Error("expected error for missing MQTT_PASS, got nil")
	}
}

func TestParseLogLevel(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  slog.Level
		isErr bool
	}{
		{"debug", slog.LevelDebug, false},
		{"info", slog.LevelInfo, false},
		{"warn", slog.LevelWarn, false},
		{"error", slog.LevelError, false},
		{"invalid", slog.LevelInfo, true},
	} {
		got, err := parseLogLevel(tc.input)
		if tc.isErr != (err != nil) {
			t.Errorf("parseLogLevel(%q): err=%v", tc.input, err)
		}
		if !tc.isErr && got != tc.want {
			t.Errorf("parseLogLevel(%q): got %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestLoadAppConfig_DefaultLogLevel(t *testing.T) {
	setEnv(t, validEnv) // no LOG_LEVEL set
	cfg, err := LoadAppConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.LogLevel.Level() != slog.LevelInfo {
		t.Errorf("expected default log level INFO, got %v", cfg.LogLevel.Level())
	}
}
