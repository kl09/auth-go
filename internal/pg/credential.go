package pg

import (
	"context"

	"github.com/jinzhu/gorm"

	auth "github.com/kl09/auth-go"
)

type CredentialRepository struct {
	*Client
}

func NewCredentialRepository(c *Client) *CredentialRepository {
	return &CredentialRepository{
		c,
	}
}

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

func (c *CredentialRepository) Create(ctx context.Context, cred *auth.Credential) error {
	return c.db.Create(cred).Error
}
