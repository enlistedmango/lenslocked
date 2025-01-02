package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Custom error types
type ErrEmailTaken struct {
	Email string
}

func (err ErrEmailTaken) Error() string {
	return fmt.Sprintf("email address %q is already taken", err.Error())
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(email, password string) (*User, error) {
	email = strings.ToLower(email)
	// Check if email is already taken
	row := us.DB.QueryRow(`
		SELECT id FROM users WHERE email = $1`, email)
	var id int
	err := row.Scan(&id)
	if err != sql.ErrNoRows {
		if err == nil {
			return nil, ErrEmailTaken{Email: email}
		}
		return nil, err
	}

	// Validate password
	if len(password) < 8 {
		return nil, fmt.Errorf("password must be at least 8 characters long")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}
	passwordHash := string(hashedBytes)

	user := User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	row = us.DB.QueryRow(`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2) RETURNING id, created_at, updated_at`,
		user.Email, user.PasswordHash)
	err = row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	return &user, nil
}

type contextKey string

const userContextKey contextKey = "user"

func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func UserFromContext(ctx context.Context) *User {
	user, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		return nil
	}
	return user
}

func (us *UserService) GetByID(id int) (*User, error) {
	user := &User{
		ID: id,
	}
	row := us.DB.QueryRow(`
		SELECT email, password_hash, created_at, updated_at
		FROM users WHERE id = $1`, id)
	err := row.Scan(
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)
	user := User{
		Email: email,
	}
	row := us.DB.QueryRow(`
		SELECT id, password_hash, created_at, updated_at 
		FROM users WHERE email = $1`, email)
	err := row.Scan(&user.ID, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}
	return &user, nil
}
