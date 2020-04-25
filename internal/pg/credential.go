package pg

import (
	"context"

	"github.com/jinzhu/gorm"

	auth "github.com/kl09/auth-go"
)

// CredentialRepository is a repository for credentials.
type CredentialRepository struct {
	*Client
}

// NewCredentialRepository creates a new CredentialRepository.
func NewCredentialRepository(c *Client) *CredentialRepository {
	return &CredentialRepository{
		c,
	}
}

// ByToken returns a Credential by token.
func (c *CredentialRepository) ByToken(ctx context.Context, token string) (auth.Credential, error) {
	cred := auth.Credential{}

	db := c.db.Where("token = ?", token).Take(&cred)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return cred, auth.NewError(auth.ErrCredNotFound, "Credential not found")
		}

		return cred, db.Error
	}

	return cred, nil
}

// ByID returns a Credential by id.
func (c *CredentialRepository) ByID(ctx context.Context, id int) (auth.Credential, error) {
	cred := auth.Credential{}

	db := c.db.Where("id = ?", id).Take(&cred)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return cred, auth.NewError(auth.ErrCredNotFound, "Credential not found")
		}

		return cred, db.Error
	}

	return cred, nil
}

// ByEmail returns a Credential by email.
func (c *CredentialRepository) ByEmail(ctx context.Context, email string) (auth.Credential, error) {
	cred := auth.Credential{}

	db := c.db.Where("email = ?", email).Take(&cred)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return cred, auth.NewError(auth.ErrCredNotFound, "Credential not found")
		}

		return cred, db.Error
	}

	return cred, nil
}

// Create creates a new Credential.
func (c *CredentialRepository) Create(ctx context.Context, cred *auth.Credential) error {
	return c.db.Create(cred).Error
}
