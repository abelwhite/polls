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
func (m *UserModel) Authenticate(email, password string) (int, error) {
	//compare
	var id int
	var hashedPassword []byte

	//check if there is a row in the table for the email provided
	query := `
		SELECT users_id, password_hash
		FROM users 
		WHERE email = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	} //handling error
	//the user does exist
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	//password is correct
	return id, nil
}
