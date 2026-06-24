package server

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"GoApp/internal/broker"
	config "GoApp/internal/config"
	"GoApp/internal/database"
	"GoApp/internal/views"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	_ "github.com/joho/godotenv/autoload"
	mqtt "github.com/mochi-mqtt/server/v2"
)

type DB interface {
	Health() map[string]string
	CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error)
	GetUserByEmail(ctx context.Context, email string) (database.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (database.User, error)
	UpdateUserName(ctx context.Context, arg database.UpdateUserNameParams) (database.User, error)
	UpdateUserPassword(ctx context.Context, arg database.UpdateUserPasswordParams) error
	DeleteUser(ctx context.Context, id uuid.UUID) error

	CreateSession(ctx context.Context, arg database.CreateSessionParams) (database.Session, error)
	GetSessionByToken(ctx context.Context, token string) (database.Session, error)
	DeleteSession(ctx context.Context, token string) error
	DeleteExpiredSessions(ctx context.Context) error
	GetActiveSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]database.Session, error)
	DeleteSessionByID(ctx context.Context, arg database.DeleteSessionByIDParams) error

	CreateContact(ctx context.Context, arg database.CreateContactParams) (database.Contact, error)
	CountContactsByIPToday(ctx context.Context, ipAddress string) (int64, error)
	CountContactsByEmailToday(ctx context.Context, email string) (int64, error)

	InsertAirTempHumidReading(ctx context.Context, arg database.InsertAirTempHumidReadingParams) error
	GetLatestAirTempHumidReadings(ctx context.Context) ([]database.GetLatestAirTempHumidReadingsRow, error)
	GetAirTempHumidReadingsByAddr(ctx context.Context, arg database.GetAirTempHumidReadingsByAddrParams) ([]database.GetAirTempHumidReadingsByAddrRow, error)
	DeleteOldAirTempHumidReadings(ctx context.Context, createdAt time.Time) error

	// MQTT credentials
	GetMQTTCredentialByUsername(ctx context.Context, username string) (database.MqttCredential, error)
	CreateMQTTCredential(ctx context.Context, arg database.CreateMQTTCredentialParams) (database.MqttCredential, error)
	GetMQTTACLByCredentialID(ctx context.Context, credentialID uuid.UUID) ([]database.MqttAcl, error)
	CreateMQTTACL(ctx context.Context, arg database.CreateMQTTACLParams) (database.MqttAcl, error)

	// Soil readings
	InsertSoilMoistureReading(ctx context.Context, arg database.InsertSoilMoistureReadingParams) error
	GetLatestSoilMoistureReadings(ctx context.Context) ([]database.GetLatestSoilMoistureReadingsRow, error)
	DeleteOldSoilMoistureReadings(ctx context.Context, createdAt time.Time) error
	GetSoilMoistureReadingsBySensorIdx(ctx context.Context, arg database.GetSoilMoistureReadingsBySensorIdxParams) ([]database.GetSoilMoistureReadingsBySensorIdxRow, error)
}
type sqlDB struct {
	raw     *sql.DB
	queries *database.Queries
}

func (s *sqlDB) Health() map[string]string {
	return database.Health(s.raw)
}

func (s *sqlDB) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	return s.queries.CreateUser(ctx, arg)
}

func (s *sqlDB) GetUserByEmail(ctx context.Context, email string) (database.User, error) {
	return s.queries.GetUserByEmail(ctx, email)
}

func (s *sqlDB) GetUserByID(ctx context.Context, id uuid.UUID) (database.User, error) {
	return s.queries.GetUserByID(ctx, id)
}

func (s *sqlDB) UpdateUserName(ctx context.Context, arg database.UpdateUserNameParams) (database.User, error) {
	return s.queries.UpdateUserName(ctx, arg)
}

func (s *sqlDB) UpdateUserPassword(ctx context.Context, arg database.UpdateUserPasswordParams) error {
	return s.queries.UpdateUserPassword(ctx, arg)
}

func (s *sqlDB) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.queries.DeleteUser(ctx, id)
}

func (s *sqlDB) CreateSession(ctx context.Context, arg database.CreateSessionParams) (database.Session, error) {
	return s.queries.CreateSession(ctx, arg)
}

func (s *sqlDB) GetSessionByToken(ctx context.Context, token string) (database.Session, error) {
	return s.queries.GetSessionByToken(ctx, token)
}

func (s *sqlDB) DeleteSession(ctx context.Context, token string) error {
	return s.queries.DeleteSession(ctx, token)
}

func (s *sqlDB) DeleteExpiredSessions(ctx context.Context) error {
	return s.queries.DeleteExpiredSessions(ctx)
}

func (s *sqlDB) GetActiveSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]database.Session, error) {
	return s.queries.GetActiveSessionsByUserID(ctx, userID)
}

func (s *sqlDB) DeleteSessionByID(ctx context.Context, arg database.DeleteSessionByIDParams) error {
	return s.queries.DeleteSessionByID(ctx, arg)
}

func (s *sqlDB) CreateContact(ctx context.Context, arg database.CreateContactParams) (database.Contact, error) {
	return s.queries.CreateContact(ctx, arg)
}

func (s *sqlDB) CountContactsByIPToday(ctx context.Context, ipAddress string) (int64, error) {
	return s.queries.CountContactsByIPToday(ctx, ipAddress)
}

func (s *sqlDB) CountContactsByEmailToday(ctx context.Context, email string) (int64, error) {
	return s.queries.CountContactsByEmailToday(ctx, email)
}

func (s *sqlDB) InsertAirTempHumidReading(ctx context.Context, arg database.InsertAirTempHumidReadingParams) error {
	return s.queries.InsertAirTempHumidReading(ctx, arg)
}

func (s *sqlDB) GetLatestAirTempHumidReadings(ctx context.Context) ([]database.GetLatestAirTempHumidReadingsRow, error) {
	return s.queries.GetLatestAirTempHumidReadings(ctx)
}

func (s *sqlDB) GetAirTempHumidReadingsByAddr(ctx context.Context, arg database.GetAirTempHumidReadingsByAddrParams) ([]database.GetAirTempHumidReadingsByAddrRow, error) {
	return s.queries.GetAirTempHumidReadingsByAddr(ctx, arg)
}

func (s *sqlDB) DeleteOldAirTempHumidReadings(ctx context.Context, createdAt time.Time) error {
	return s.queries.DeleteOldAirTempHumidReadings(ctx, createdAt)
}

func (s *sqlDB) GetMQTTCredentialByUsername(ctx context.Context, username string) (database.MqttCredential, error) {
	return s.queries.GetMQTTCredentialByUsername(ctx, username)
}
func (s *sqlDB) CreateMQTTCredential(ctx context.Context, arg database.CreateMQTTCredentialParams) (database.MqttCredential, error) {
	return s.queries.CreateMQTTCredential(ctx, arg)
}
func (s *sqlDB) GetMQTTACLByCredentialID(ctx context.Context, credentialID uuid.UUID) ([]database.MqttAcl, error) {
	return s.queries.GetMQTTACLByCredentialID(ctx, credentialID)
}
func (s *sqlDB) CreateMQTTACL(ctx context.Context, arg database.CreateMQTTACLParams) (database.MqttAcl, error) {
	return s.queries.CreateMQTTACL(ctx, arg)
}
func (s *sqlDB) InsertSoilMoistureReading(ctx context.Context, arg database.InsertSoilMoistureReadingParams) error {
	return s.queries.InsertSoilMoistureReading(ctx, arg)
}
func (s *sqlDB) GetLatestSoilMoistureReadings(ctx context.Context) ([]database.GetLatestSoilMoistureReadingsRow, error) {
	return s.queries.GetLatestSoilMoistureReadings(ctx)
}
func (s *sqlDB) DeleteOldSoilMoistureReadings(ctx context.Context, createdAt time.Time) error {
	return s.queries.DeleteOldSoilMoistureReadings(ctx, createdAt)
}
func (s *sqlDB) GetSoilMoistureReadingsBySensorIdx(ctx context.Context, arg database.GetSoilMoistureReadingsBySensorIdxParams) ([]database.GetSoilMoistureReadingsBySensorIdxRow, error) {
	return s.queries.GetSoilMoistureReadingsBySensorIdx(ctx, arg)
}

type Server struct {
	port       int
	db         DB
	cfg        *config.Config
	hub        *Hub
	wsUpgrader websocket.Upgrader
}

func NewServer(cfg *config.Config) (*http.Server, *mqtt.Server, error) {
	dbCfg := &database.DBConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		Database: cfg.DBDatabase,
		Username: cfg.DBUsername,
		Password: cfg.DBPassword,
		Schema:   cfg.DBSchema,
	}

	raw := database.Open(dbCfg)

	if err := database.Migrate(raw); err != nil {
		return nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	s := &Server{
		port: cfg.Port,
		db: &sqlDB{
			raw:     raw,
			queries: database.New(raw),
		},
		cfg: cfg,
		hub: NewHub(cfg),
	}
	s.wsUpgrader = websocket.Upgrader{
		HandshakeTimeout: wsHandshakeTimeout,
		CheckOrigin:      s.wsCheckOrigin,
	}
	s.StartSessionCleanup(context.Background(), 1*time.Hour)

	var mqttTLS *tls.Config

	// Use Let's Encrypt overrides if provided; otherwise fallback to internal mTLS certs
	mqttCertPath := cfg.MQTTCertPath
	mqttKeyPath := cfg.MQTTKeyPath
	if mqttCertPath == "" && cfg.AppEnv == config.EnvProduction {
		// prod only: fall back to mTLS certs if no dedicated MQTT cert
		mqttCertPath = cfg.TLSCertPath
		mqttKeyPath = cfg.TLSKeyPath
	}

	if mqttCertPath != "" {
		cert, err := tls.LoadX509KeyPair(mqttCertPath, mqttKeyPath)
		if err != nil {
			return nil, nil, fmt.Errorf("mqtt tls: load cert: %w", err)
		}
		mqttTLS = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}
	}

	mqttSrv, err := broker.Start(cfg.MQTTTLSPort, cfg.MQTTPort, s.db, s.hub, mqttTLS, cfg.MQTTUser, cfg.MQTTPass)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to start MQTT broker: %w", err)
	}

	httpSrv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.TLSPort),
		Handler:      s.RegisterRoutes(cfg),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return httpSrv, mqttSrv, nil
}

func (s *Server) siteConfig(c *gin.Context) views.SiteConfig {
	return views.SiteConfig{
		AppName:      s.cfg.AppName,
		ContactEmail: s.cfg.ContactEmail,
		ContactPhone: s.cfg.ContactPhone,
		BaseURL:      s.requestBaseURL(c),
	}
}
