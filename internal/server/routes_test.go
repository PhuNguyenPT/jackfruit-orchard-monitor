package server

import (
	"GoApp/internal/database"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type mockDB struct {
	deleteExpiredSessionsCalled int
}

func (m *mockDB) Health() map[string]string {
	return map[string]string{"status": "up", "message": "It's healthy"}
}

func (m *mockDB) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return database.User{}, err
	}
	return database.User{ID: id, Name: arg.Name, Email: arg.Email, PasswordHash: arg.PasswordHash, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (m *mockDB) GetUserByEmail(ctx context.Context, email string) (database.User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return database.User{}, err
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	return database.User{
		ID:           id,
		Name:         "Test User",
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (m *mockDB) GetUserByID(ctx context.Context, id uuid.UUID) (database.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	return database.User{
		ID:           id,
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (m *mockDB) UpdateUserName(ctx context.Context, arg database.UpdateUserNameParams) (database.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return database.User{}, err
	}
	return database.User{
		ID:           arg.ID,
		Name:         arg.Name,
		Email:        "test@example.com",
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (m *mockDB) UpdateUserPassword(ctx context.Context, arg database.UpdateUserPasswordParams) error {
	return nil
}

func (m *mockDB) DeleteUser(ctx context.Context, id uuid.UUID) error { return nil }

func (m *mockDB) CreateSession(ctx context.Context, arg database.CreateSessionParams) (database.Session, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return database.Session{}, err
	}
	return database.Session{ID: id, UserID: arg.UserID, Token: arg.Token, CreatedAt: time.Now(), ExpiresAt: time.Now().Add(1 * time.Hour.Abs())}, nil
}

func (m *mockDB) GetSessionByToken(ctx context.Context, token string) (database.Session, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return database.Session{}, err
	}
	userID, err := uuid.NewV7()
	if err != nil {
		return database.Session{}, err
	}
	return database.Session{ID: id, UserID: userID, Token: token, CreatedAt: time.Now(), ExpiresAt: time.Now().Add(1 * time.Hour.Abs())}, nil
}

func (m *mockDB) DeleteSession(ctx context.Context, token string) error { return nil }

func (m *mockDB) DeleteExpiredSessions(ctx context.Context) error {
	m.deleteExpiredSessionsCalled++
	return nil
}

func (m *mockDB) GetActiveSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]database.Session, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return []database.Session{}, err
	}
	token, err := uuid.NewV7()
	if err != nil {
		return []database.Session{}, err
	}
	return []database.Session{
		{ID: id, UserID: userID, Token: token.String(), UserAgent: "Mozilla/5.0 Test Browser", CreatedAt: time.Now(), ExpiresAt: time.Now().Add(1 * time.Hour.Abs()), IpAddress: "127.0.0.1"},
	}, nil
}

func (m *mockDB) DeleteSessionByID(ctx context.Context, arg database.DeleteSessionByIDParams) error {
	return nil
}

func (m *mockDB) CreateContact(ctx context.Context, arg database.CreateContactParams) (database.Contact, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return database.Contact{}, err
	}
	return database.Contact{
		ID:        id,
		Name:      arg.Name,
		Email:     arg.Email,
		Subject:   arg.Subject,
		Message:   arg.Message,
		IpAddress: "127.0.0.1",
		CreatedAt: time.Now(),
	}, nil
}

func (m *mockDB) CountContactsByIPToday(ctx context.Context, ipAddress string) (int64, error) {
	return 0, nil
}

func (m *mockDB) CountContactsByEmailToday(ctx context.Context, email string) (int64, error) {
	return 0, nil
}

func (m *mockDB) InsertSensorReading(ctx context.Context, arg database.InsertSensorReadingParams) error {
	return nil
}

func (m *mockDB) GetLatestReadings(ctx context.Context) ([]database.GetLatestReadingsRow, error) {
	return []database.GetLatestReadingsRow{
		{Addr: 1, Temperature: 284, Humidity: 742, CreatedAt: time.Now()},
		{Addr: 2, Temperature: 291, Humidity: 718, CreatedAt: time.Now()},
		{Addr: 3, Temperature: 276, Humidity: 765, CreatedAt: time.Now()},
	}, nil
}

func (m *mockDB) GetReadingsByAddr(ctx context.Context, arg database.GetReadingsByAddrParams) ([]database.GetReadingsByAddrRow, error) {
	return []database.GetReadingsByAddrRow{
		{Addr: arg.Addr, Temperature: 284, Humidity: 742, CreatedAt: time.Now()},
		{Addr: arg.Addr, Temperature: 281, Humidity: 748, CreatedAt: time.Now().Add(-1 * time.Minute)},
		{Addr: arg.Addr, Temperature: 279, Humidity: 751, CreatedAt: time.Now().Add(-2 * time.Minute)},
	}, nil
}

func (m *mockDB) DeleteOldSensorReadings(ctx context.Context, createdAt time.Time) error {
	return nil
}

// ---------------------------------------------------------------------------

var testHandler http.Handler

func TestMain(m *testing.M) {
	if err := os.Chdir("../../"); err != nil {
		log.Fatalf("failed to change directory: %v", err)
	}
	s := &Server{
		db:  &mockDB{},
		cfg: &Config{AppEnv: EnvTest, GinMode: gin.TestMode},
	}
	testHandler = s.RegisterRoutes(s.cfg)
	os.Exit(m.Run())
}

func TestHomePageHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	testHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Errorf("expected HTML content type, got %v", ct)
	}
}

func TestUnknownRoute(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/does-not-exist", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	testHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %v", rr.Code)
	}
}
