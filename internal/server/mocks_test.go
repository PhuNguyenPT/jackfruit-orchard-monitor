package server

import (
	"GoApp/internal/database"
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type mockDB struct {
	deleteExpiredSessionsCalled int
}

// --- Core Server & User Mocks ---
func (m *mockDB) Health() map[string]string {
	return map[string]string{"status": "up", "message": "It's healthy"}
}

func (m *mockDB) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	return database.User{
		ID:           uuid.Must(uuid.NewV7()),
		Name:         arg.Name,
		Email:        arg.Email,
		PasswordHash: arg.PasswordHash,
	}, nil
}

func (m *mockDB) GetUserByEmail(ctx context.Context, email string) (database.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return database.User{}, err
	}

	return database.User{
		ID:           uuid.Must(uuid.NewV7()),
		Name:         "Test User",
		Email:        email,
		PasswordHash: string(hash),
	}, nil
}

func (m *mockDB) GetUserByID(ctx context.Context, id uuid.UUID) (database.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return database.User{}, err
	}

	return database.User{
		ID:           id,
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: string(hash),
	}, nil
}

func (m *mockDB) UpdateUserName(ctx context.Context, arg database.UpdateUserNameParams) (database.User, error) {
	return database.User{ID: arg.ID, Name: arg.Name, Email: "test@example.com"}, nil
}

func (m *mockDB) UpdateUserPassword(ctx context.Context, arg database.UpdateUserPasswordParams) error {
	return nil
}
func (m *mockDB) DeleteUser(ctx context.Context, id uuid.UUID) error { return nil }

// --- Session Mocks ---
func (m *mockDB) CreateSession(ctx context.Context, arg database.CreateSessionParams) (database.Session, error) {
	return database.Session{
		ID:        uuid.Must(uuid.NewV7()),
		UserID:    arg.UserID,
		Token:     arg.Token,
		ExpiresAt: time.Now().Add(time.Hour),
	}, nil
}

func (m *mockDB) GetSessionByToken(ctx context.Context, token string) (database.Session, error) {
	return database.Session{
		ID:        uuid.Must(uuid.NewV7()),
		UserID:    uuid.Must(uuid.NewV7()),
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour),
	}, nil
}

func (m *mockDB) DeleteSession(ctx context.Context, token string) error { return nil }
func (m *mockDB) DeleteExpiredSessions(ctx context.Context) error {
	m.deleteExpiredSessionsCalled++
	return nil
}
func (m *mockDB) GetActiveSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]database.Session, error) {
	return []database.Session{
		{
			ID:        uuid.Must(uuid.NewV7()),
			UserID:    userID,
			Token:     "mock-token",
			IpAddress: "127.0.0.1",
			UserAgent: "Mozilla/5.0 Test Browser",
		},
	}, nil
}
func (m *mockDB) DeleteSessionByID(ctx context.Context, arg database.DeleteSessionByIDParams) error {
	return nil
}

// --- Contact Mocks ---
func (m *mockDB) CreateContact(ctx context.Context, arg database.CreateContactParams) (database.Contact, error) {
	return database.Contact{
		ID:    uuid.Must(uuid.NewV7()),
		Name:  arg.Name,
		Email: arg.Email,
	}, nil
}
func (m *mockDB) CountContactsByIPToday(ctx context.Context, ipAddress string) (int64, error) {
	return 0, nil
}
func (m *mockDB) CountContactsByEmailToday(ctx context.Context, email string) (int64, error) {
	return 0, nil
}

// --- Sensor & Soil Reading Mocks ---
func (m *mockDB) InsertAirTempHumidReading(ctx context.Context, arg database.InsertAirTempHumidReadingParams) error {
	return nil
}

func (m *mockDB) GetLatestAirTempHumidReadings(ctx context.Context) ([]database.GetLatestAirTempHumidReadingsRow, error) {
	return []database.GetLatestAirTempHumidReadingsRow{
		{Addr: 1, Temperature: 284, Humidity: 742, CreatedAt: time.Now()},
		{Addr: 2, Temperature: 291, Humidity: 718, CreatedAt: time.Now()},
		{Addr: 3, Temperature: 276, Humidity: 765, CreatedAt: time.Now()},
	}, nil
}

func (m *mockDB) GetAirTempHumidReadingsByAddr(ctx context.Context, arg database.GetAirTempHumidReadingsByAddrParams) ([]database.GetAirTempHumidReadingsByAddrRow, error) {
	return []database.GetAirTempHumidReadingsByAddrRow{
		{Addr: arg.Addr, Temperature: 284, Humidity: 742, CreatedAt: time.Now()},
		{Addr: arg.Addr, Temperature: 281, Humidity: 748, CreatedAt: time.Now().Add(-1 * time.Minute)},
		{Addr: arg.Addr, Temperature: 279, Humidity: 751, CreatedAt: time.Now().Add(-2 * time.Minute)},
	}, nil
}

func (m *mockDB) DeleteOldAirTempHumidReadings(ctx context.Context, createdAt time.Time) error {
	return nil
}

func (m *mockDB) InsertSoilMoistureReading(ctx context.Context, arg database.InsertSoilMoistureReadingParams) error {
	return nil
}

func (m *mockDB) GetLatestSoilMoistureReadings(ctx context.Context) ([]database.GetLatestSoilMoistureReadingsRow, error) {
	return []database.GetLatestSoilMoistureReadingsRow{
		{SensorIdx: 0, Raw: 1500, CreatedAt: time.Now()},
		{SensorIdx: 1, Raw: 3000, CreatedAt: time.Now()},
	}, nil
}

func (m *mockDB) DeleteOldSoilMoistureReadings(ctx context.Context, createdAt time.Time) error {
	return nil
}

// --- MQTT Mocks ---
func (m *mockDB) GetMQTTCredentialByUsername(ctx context.Context, username string) (database.MqttCredential, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte("mqtt_password123"), bcrypt.DefaultCost)
	if err != nil {
		return database.MqttCredential{}, err
	}

	return database.MqttCredential{
		ID:        uuid.Must(uuid.NewV7()),
		Username:  username,
		Password:  string(hash),
		CreatedAt: time.Now(),
	}, nil
}

func (m *mockDB) CreateMQTTCredential(ctx context.Context, arg database.CreateMQTTCredentialParams) (database.MqttCredential, error) {
	return database.MqttCredential{
		ID:        uuid.Must(uuid.NewV7()),
		Username:  arg.Username,
		Password:  arg.Password,
		CreatedAt: time.Now(),
	}, nil
}

func (m *mockDB) GetMQTTACLByCredentialID(ctx context.Context, credentialID uuid.UUID) ([]database.MqttAcl, error) {
	return []database.MqttAcl{
		{ID: uuid.Must(uuid.NewV7()), CredentialID: credentialID, Topic: "sht40/+/data", Permission: "w"},
		{ID: uuid.Must(uuid.NewV7()), CredentialID: credentialID, Topic: "soil/#", Permission: "r"},
	}, nil
}

func (m *mockDB) CreateMQTTACL(ctx context.Context, arg database.CreateMQTTACLParams) (database.MqttAcl, error) {
	perm := arg.Permission
	if perm == "" {
		perm = "r"
	}

	return database.MqttAcl{
		ID:           uuid.Must(uuid.NewV7()),
		CredentialID: arg.CredentialID,
		Topic:        arg.Topic,
		Permission:   perm,
	}, nil
}
