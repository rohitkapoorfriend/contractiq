package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/contractiq/contractiq/pkg/apperror"
	"github.com/contractiq/contractiq/pkg/identifier"
	"github.com/jmoiron/sqlx"
)

// User represents an authenticated user.
type User struct {
	ID           string    `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	Name         string    `db:"name" json:"name"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// RegisterRequest is the input for user registration.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// LoginRequest is the input for user login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse is returned after successful login/registration.
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// UserService handles user authentication.
type UserService struct {
	db  *sqlx.DB
	jwt *JWTService
}

// NewUserService creates a new user service.
func NewUserService(db *sqlx.DB, jwt *JWTService) *UserService {
	return &UserService{db: db, jwt: jwt}
}

// Register creates a new user account.
func (s *UserService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	// Check if email exists
	var exists bool
	err := s.db.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, apperror.NewConflict("email already registered")
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := User{
		ID:           identifier.New(),
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: hash,
		CreatedAt:    time.Now().UTC(),
	}

	_, err = s.db.ExecContext(ctx,
		"INSERT INTO users (id, email, name, password_hash, created_at) VALUES ($1, $2, $3, $4, $5)",
		user.ID, user.Email, user.Name, user.PasswordHash, user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := s.jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: user}, nil
}

// Login authenticates a user and returns a JWT token.
func (s *UserService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	var user User
	err := s.db.GetContext(ctx, &user,
		"SELECT id, email, name, password_hash, created_at FROM users WHERE email = $1", req.Email,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewUnauthorized("invalid email or password")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if !CheckPassword(req.Password, user.PasswordHash) {
		return nil, apperror.NewUnauthorized("invalid email or password")
	}

	token, err := s.jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: user}, nil
}
