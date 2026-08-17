package server

import (
	"context"
	"testing"

	appConfig "GoApp/internal/config"
)

func TestSeedDefaultAdmin_CreatesWhenMissing(t *testing.T) {
	db := &mockDB{missingEmails: map[string]bool{"admin@example.com": true}}
	s := &Server{db: db, cfg: &appConfig.Config{
		AdminEmail:    "admin@example.com",
		AdminPassword: "supersecret1",
		AdminName:     "Admin",
	}}

	if err := s.seedDefaultAdmin(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(db.addedToGroup) != 1 {
		t.Fatalf("expected 1 group assignment, got %d", len(db.addedToGroup))
	}
}

func TestSeedDefaultAdmin_ExistingUser_StillEnsuresGroupMembership(t *testing.T) {
	db := &mockDB{}
	s := &Server{db: db, cfg: &appConfig.Config{
		AdminEmail:    "admin@example.com",
		AdminPassword: "supersecret1",
		AdminName:     "Admin",
	}}

	if err := s.seedDefaultAdmin(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(db.addedToGroup) != 1 {
		t.Fatalf("expected group assignment to be ensured for existing user, got %d calls", len(db.addedToGroup))
	}
}

func TestSeedDefaultAdmin_NoopWhenNotConfigured(t *testing.T) {
	db := &mockDB{}
	s := &Server{db: db, cfg: &appConfig.Config{}} // AdminEmail unset

	if err := s.seedDefaultAdmin(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(db.addedToGroup) != 0 {
		t.Errorf("expected no seeding when AdminEmail unset")
	}
}
