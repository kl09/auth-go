package mock

import (
	"context"

	auth "github.com/kl09/auth-go"
)

type CredentialRepository struct {
	ByTokenFn func(ctx context.Context, token string) (auth.Credential, error)
	ByIDFn    func(ctx context.Context, id int) (auth.Credential, error)
	ByEmailFn func(ctx context.Context, email string) (auth.Credential, error)
	CreateFn  func(ctx context.Context, c *auth.Credential) error
}

func (c *CredentialRepository) ByToken(ctx context.Context, token string) (auth.Credential, error) {
	return c.ByTokenFn(ctx, token)
}

func (c *CredentialRepository) ByID(ctx context.Context, id int) (auth.Credential, error) {
	return c.ByIDFn(ctx, id)
}

func (c *CredentialRepository) ByEmail(ctx context.Context, email string) (auth.Credential, error) {
	return c.ByEmailFn(ctx, email)
}

func (c *CredentialRepository) Create(ctx context.Context, cred *auth.Credential) error {
	return c.CreateFn(ctx, cred)
}
