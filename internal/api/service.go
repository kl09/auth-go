package api

import (
	"context"
	"time"

	auth "github.com/kl09/auth-go"
)

const tokenLength = 128

type CredentialService struct {
	credentialRepository auth.CredentialRepository
	nowFn                func() time.Time
	generatorFn          func(n int) (string, error)
}

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

func (c *CredentialService) ByToken(ctx context.Context, token string) (auth.Credential, error) {
	return c.credentialRepository.ByToken(ctx, token)
}

func (c *CredentialService) Register(ctx context.Context, cred *auth.Credential) error {
	var err error

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
