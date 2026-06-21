package server

import (
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
	Port    int    `validate:"required,gt=0"`
	AppEnv  string `validate:"required,oneof=dev production test"`
	AppName string `validate:"required"`
	GinMode string `validate:"required,oneof=debug release test"`
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
	MQTTPort     int    `validate:"required,gt=0"`
	MQTTUser     string `validate:"required"`
	MQTTPass     string `validate:"required,min=8"`
	MQTTCertPath string
	MQTTKeyPath  string
	// Soil Calibration Thresholds
	SoilDryValue int `validate:"required,gt=0"`
	SoilWetValue int `validate:"required,gt=0"`
	// App URL
	BaseURL string `validate:"required,url"`
	// Contact Config
	ContactEmail string `validate:"required,email"`
	ContactPhone string
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
	cfg := &Config{
		Port:         getEnvOrDefaultInt("PORT", 8080),
		AppEnv:       os.Getenv("APP_ENV"),
		AppName:      getEnvOrDefaultStr("APP_NAME", "Prizm"),
		BaseURL:      getEnvOrDefaultStr("BASE_URL", "https://yourdomain.com"),
		GinMode:      os.Getenv("GIN_MODE"),
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
		MQTTPort:     getEnvOrDefaultInt("MQTT_PORT", 8883),
		MQTTUser:     getEnvOrDefaultStr("MQTT_USER", "esp32"),
		MQTTPass:     os.Getenv("MQTT_PASS"),
		MQTTCertPath: getEnvOrDefaultStr("MQTT_CERT_PATH", ""),
		MQTTKeyPath:  getEnvOrDefaultStr("MQTT_KEY_PATH", ""),
		SoilDryValue: getEnvOrDefaultInt("SOIL_DRY_VALUE", 3340),
		SoilWetValue: getEnvOrDefaultInt("SOIL_WET_VALUE", 1805),
		ContactEmail: getEnvOrDefaultStr("CONTACT_EMAIL", "info@yourdomain.com"),
		ContactPhone: getEnvOrDefaultStr("CONTACT_PHONE", ""),
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
