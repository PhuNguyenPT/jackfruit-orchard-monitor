package server

import (
	"GoApp/internal/database"
	"GoApp/internal/rbac"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) seedDefaultAdmin(ctx context.Context) error {
	if s.cfg.AdminEmail == "" {
		return nil
	}

	group, err := s.db.GetGroupByName(ctx, string(rbac.GroupAdmin))
	if err != nil {
		return fmt.Errorf("seed default admin: lookup admin group: %w", err)
	}

	user, err := s.db.GetUserByEmail(ctx, s.cfg.AdminEmail)
	switch {
	case err == nil:
		if err := s.db.AddUserToGroup(ctx, database.AddUserToGroupParams{UserID: user.ID, GroupID: group.ID}); err != nil {
			return fmt.Errorf("seed default admin: assign group: %w", err)
		}
		return nil
	case errors.Is(err, sql.ErrNoRows):
		// fall through to create
	default:
		return fmt.Errorf("seed default admin: lookup existing user: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(s.cfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("seed default admin: hash password: %w", err)
	}

	user, err = s.db.CreateUser(ctx, database.CreateUserParams{
		Email:        s.cfg.AdminEmail,
		PasswordHash: string(hash),
		Name:         s.cfg.AdminName,
	})
	if err != nil {
		return fmt.Errorf("seed default admin: create user: %w", err)
	}

	if err := s.db.AddUserToGroup(ctx, database.AddUserToGroupParams{UserID: user.ID, GroupID: group.ID}); err != nil {
		return fmt.Errorf("seed default admin: assign group: %w", err)
	}

	log.Printf("[Server] seeded default admin account: %s", s.cfg.AdminEmail)
	return nil
}
