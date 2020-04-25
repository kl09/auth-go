package api

import (
	"context"
	"time"

	auth "github.com/kl09/auth-go"
)

const tokenLength = 128

// CredentialService is a service that works with credentials.
type CredentialService struct {
	credentialRepository auth.CredentialRepository
	nowFn                func() time.Time
	generatorFn          func(n int) (string, error)
}

// NewCredentialService creates a CredentialService.
func NewCredentialService(
	r auth.CredentialRepository,
	nowFn func() time.Time,
	generatorFn func(n int) (string, error),
) *CredentialService {
	return &CredentialService{
		credentialRepository: r,
		nowFn:                nowFn,
		generatorFn:          generatorFn,
	}
}

// ByToken retrieves a Credential by token.
func (c *CredentialService) ByToken(ctx context.Context, token string) (auth.Credential, error) {
	return c.credentialRepository.ByToken(ctx, token)
}

// Register creates a new credential.
func (c *CredentialService) Register(ctx context.Context, cred *auth.Credential) error {
	var err error

	_, err = c.credentialRepository.ByEmail(ctx, cred.Email)
	if err == nil {
		return auth.NewError(auth.ErrEmailExists, "User with this email already exists.")
	}

	if auth.ErrorCode(err) != auth.ErrCredNotFound {
		return auth.WrapError(err, auth.ErrInternal, "Register failed")
	}

	cred.Password, err = hashAndSalt(cred.Password)
	if err != nil {
		return err
	}

	cred.Token, err = c.generatorFn(tokenLength)
	if err != nil {
		return err
	}

	cred.CreatedAt = c.nowFn()
	cred.UpdatedAt = c.nowFn()

	return c.credentialRepository.Create(ctx, cred)
}

// Auth checks user's email/pass.
func (c *CredentialService) Auth(ctx context.Context, email, plainPassword string) (auth.Credential, error) {
	cred, err := c.credentialRepository.ByEmail(ctx, email)
	if err != nil {
		return auth.Credential{}, auth.WrapError(err, auth.ErrAuth, "Auth failed")
	}

	result := comparePasswords(cred.Password, plainPassword)
	if !result {
		return auth.Credential{}, auth.NewError(auth.ErrAuth, "Auth failed")
	}

	return cred, nil
}
