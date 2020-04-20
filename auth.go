package auth

import (
	"context"
	"time"
)

// Credential is a user's credential.
type Credential struct {
	ID                       int
	Password                 string
	Token                    string
	Email                    string
	EmailTmp                 string
	EmailVerified            bool
	VerificationCode         string
	VerificationCodeAttempts uint8
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

// CredentialRepository is a storage for credentials.
type CredentialRepository interface {
	// ByToken retrieves a Credential by token.
	ByToken(ctx context.Context, token string) (Credential, error)
	// ByID retrieves a Credential by id.
	ByID(ctx context.Context, id int) (Credential, error)
	// ByEmail retrieves a Credential by email.
	ByEmail(ctx context.Context, email string) (Credential, error)
	// Create creates a new Credential without verification.
	Create(ctx context.Context, c *Credential) error
}

// CredentialService represents a service for credentials.
type CredentialService interface {
	// ByToken retrieves a Credential by token.
	ByToken(ctx context.Context, token string) (Credential, error)
	// Register creates a new Credential without verification.
	Register(ctx context.Context, c *Credential) error
	// Auth makes an auth attempt.
	Auth(ctx context.Context, email, plainPassword string) (Credential, error)
}
