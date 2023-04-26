// Filename: internal/models/users.go
package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNoRecord           = errors.New("No matching record found")
	ErrInvalidCredentials = errors.New("invalid Credentials")
	ErrDuplicateEmail     = errors.New("duplicate email")
)

// Create a user
type User struct {
	UserID    int64
	Name      string
	Email     string
	Password  []byte
	Activated bool
	CreatedAt time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	//lets first hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	query := `
			INSERT INTO users(username, email, password_hash)
			VALUES($1, $2, $3)
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = m.DB.ExecContext(ctx, query, name, email, hashedPassword)
	if err != nil {
		switch {
		case err.Error() == `pgx: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}
