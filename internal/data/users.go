package data

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("Email Already Exists")
	ErrUserNotFound   = errors.New("User Not Found")
)

type User struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"-"`
	Activated    bool      `json:"activated"`
	Version      int       `json:"-"`
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) (*User, error) {

	query := `
		INSERT INTO users (name,email,password_hash,activated)
		VALUES ($1,$2,$3,$4)
		RETURNING id,created_at,version
			 `

	args := []any{user.Name, user.Email, user.PasswordHash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, user.CreatedAt, user.Version)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, ErrDuplicateEmail
		}
		return nil, err
	}
	return user, nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {

	query := `
		SELECT id, created_at, name, email, password_hash, activated, version
		FROM users
		WHERE email = $1
			 `
	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
