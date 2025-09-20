package data

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/vj-2303/gist-go/internal/validator"
)

var (
	ErrDuplicateEmail = errors.New("Email Already Exists")
	ErrUserNotFound   = errors.New("User Not Found")
	AnonymousUser     = &User{}
)

type User struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	PasswordHash []byte    `json:"-"`
	Activated    bool      `json:"activated"`
	Version      int       `json:"-"`
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) >= 3, "name", "must be atleast 3 character long")
	v.Check(len(user.Name) <= 100, "name", "cant be more than 100 chars Long")

	v.Check(user.Email != "", "email", "must be provided")
	v.Check(validator.Matches(user.Email, validator.EmailRX), "email", "must be an valid email address")

	v.Check(user.Password != "", "password", "must be provided")
	v.Check(len(user.Password) <= 72, "password", "must be less than 72 character Long")
}

func ValidateLoginUser(v *validator.Validator, email, password string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be an valid email address")

	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) <= 72, "password", "must be less than 72 character Long")
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {

	query := `
		INSERT INTO users (name,email,password_hash,activated)
		VALUES ($1,$2,$3,$4)
		RETURNING id,created_at,version
			 `

	args := []any{user.Name, user.Email, user.PasswordHash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
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

func (m UserModel) GetByID(userID int64) (*User, error) {

	query := `
		SELECT id, created_at, name, email, password_hash, activated, version
		FROM users
		WHERE id = $1
			 `
	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, userID).Scan(
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
