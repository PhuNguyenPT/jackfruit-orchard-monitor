package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
)

const (
	EnvDev        = "dev"
	EnvProduction = "production"
	EnvTest       = "test"
)

type Config struct {
	Port       int    `validate:"required,gt=0"`
	AppEnv     string `validate:"required,oneof=dev production test"`
	AppName    string `validate:"required"`
	AppVersion string
	BuildDate  string
	GinMode    string `validate:"required,oneof=debug release test"`
	LogLevel   *slog.LevelVar
	// Database
	DBHost     string `validate:"required"`
	DBPort     int    `validate:"required,gt=0"`
	DBDatabase string `validate:"required"`
	DBUsername string `validate:"required"`
	DBPassword string `validate:"required,min=8"`
	DBSchema   string `validate:"required"`
	// Pagination
	PageSize    int `validate:"gt=0"`
	MaxPageSize int `validate:"gt=0"`
	// mTLS cert paths (optional — defaults to Docker secrets paths)
	TLSCertPath string
	TLSKeyPath  string
	TLSCAPath   string
	TLSPort     int `validate:"required,gt=0"`
	// MQTT
	MQTTTLSPort  int    `validate:"required,gt=0"`
	MQTTPort     int    // plain TCP; 0 = disabled (prod)
	MQTTUser     string `validate:"required"`
	MQTTPass     string `validate:"required,min=8"`
	MQTTCertPath string
	MQTTKeyPath  string
	// Soil Calibration Thresholds
	SoilDryValue int `validate:"required,gt=0"`
	SoilWetValue int `validate:"required,gt=0"`
	// BaseURLs holds every domain this app is reachable at. The first
	// entry is the canonical/primary domain (used for things like
	// generating absolute links, sitemap, robots.txt). All entries are
	// valid WebSocket origins.
	BaseURLs []string `validate:"required,min=1,dive,url"`
	// Contact Config
	ContactEmail string `validate:"required,email"`
	ContactPhone string
}

func parseLogLevel(s string) (slog.Level, error) {
	var l slog.Level
	if err := l.UnmarshalText([]byte(s)); err != nil {
		return slog.LevelInfo, fmt.Errorf("LOG_LEVEL %q: must be debug|info|warn|error", s)
	}
	return l, nil
}

func getEnvOrDefaultJSONList(key string, defaultVal []string) ([]string, error) {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal, nil
	}
	var out []string
	if err := json.Unmarshal([]byte(v), &out); err != nil {
		return nil, fmt.Errorf("%s: invalid JSON list: %w", key, err)
	}
	return out, nil
}

func getEnvOrDefaultStr(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvOrDefaultInt(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal // not set → use default
	}
	if n, err := strconv.Atoi(v); err == nil && n > 0 {
		return n
	}
	return 0 // set but unparseable → 0, let validator catch it
}

func LoadAppConfig() (*Config, error) {
	baseURLs, err := getEnvOrDefaultJSONList("BASE_URLS", []string{"https://yourdomain.com"})
	if err != nil {
		return nil, err
	}
	logLevel, err := parseLogLevel(getEnvOrDefaultStr("LOG_LEVEL", "info"))
	if err != nil {
		return nil, err
	}
	lv := &slog.LevelVar{}
	lv.Set(logLevel)
	cfg := &Config{
		Port:         getEnvOrDefaultInt("PORT", 8080),
		AppEnv:       os.Getenv("APP_ENV"),
		AppName:      getEnvOrDefaultStr("APP_NAME", "Prizm"),
		BaseURLs:     baseURLs,
		GinMode:      os.Getenv("GIN_MODE"),
		LogLevel:     lv,
		DBHost:       os.Getenv("POSTGRES_HOST"),
		DBPort:       getEnvOrDefaultInt("POSTGRES_PORT", 5432),
		DBDatabase:   os.Getenv("POSTGRES_DATABASE"),
		DBUsername:   os.Getenv("POSTGRES_USERNAME"),
		DBPassword:   os.Getenv("POSTGRES_PASSWORD"),
		DBSchema:     os.Getenv("POSTGRES_SCHEMA"),
		PageSize:     getEnvOrDefaultInt("PAGE_SIZE", 20),
		MaxPageSize:  getEnvOrDefaultInt("MAX_PAGE_SIZE", 100),
		TLSCertPath:  getEnvOrDefaultStr("TLS_CERT_PATH", "/run/secrets/go_crt"),
		TLSKeyPath:   getEnvOrDefaultStr("TLS_KEY_PATH", "/run/secrets/go_key"),
		TLSCAPath:    getEnvOrDefaultStr("TLS_CA_PATH", "/run/secrets/backend_ca"),
		TLSPort:      getEnvOrDefaultInt("TLS_PORT", 8443),
		MQTTTLSPort:  getEnvOrDefaultInt("MQTT_TLS_PORT", 8883),
		MQTTPort:     getEnvOrDefaultInt("MQTT_PORT", 0), // 0 = disabled by default
		MQTTUser:     getEnvOrDefaultStr("MQTT_USER", "esp32"),
		MQTTPass:     os.Getenv("MQTT_PASS"),
		MQTTCertPath: getEnvOrDefaultStr("MQTT_CERT_PATH", ""),
		MQTTKeyPath:  getEnvOrDefaultStr("MQTT_KEY_PATH", ""),
		SoilDryValue: getEnvOrDefaultInt("SOIL_DRY_VALUE", 3500),
		SoilWetValue: getEnvOrDefaultInt("SOIL_WET_VALUE", 1760),
		ContactEmail: getEnvOrDefaultStr("CONTACT_EMAIL", "info@yourdomain.com"),
		ContactPhone: getEnvOrDefaultStr("CONTACT_PHONE", ""),
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
