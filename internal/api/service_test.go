package api

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	auth "github.com/kl09/auth-go"
	"github.com/kl09/auth-go/internal/mock"
)

func TestCredentialService_Register(t *testing.T) {
	cred := auth.Credential{
		Password: "12345",
		Email:    "example@example.org",
	}

	plainPass := cred.Password

	s := NewCredentialService(&mock.CredentialRepository{
		CreateFn: func(ctx context.Context, c *auth.Credential) error {
			c.ID = 1
			return nil
		},
		ByEmailFn: func(ctx context.Context, email string) (auth.Credential, error) {
			return auth.Credential{}, auth.NewError(auth.ErrCredNotFound, "Credential not found")
		},
	},
		nowFunc,
		func(n int) (string, error) {
			return "1234abcd", nil
		},
	)
	err := s.Register(context.Background(), &cred)
	if err != nil {
		t.Fatal(err)
	}

	require.NotNil(t, cred.Password)
	require.NotNil(t, cred.Token)
	require.Equal(t, cred.Token, "1234abcd")

	require.True(t, comparePasswords(cred.Password, plainPass))

	require.Equal(t, now.String(), cred.CreatedAt.String())
	require.Equal(t, now.String(), cred.UpdatedAt.String())
}

func TestCredentialService_Auth(t *testing.T) {
	token := "1234abcd"

	hash, err := hashAndSalt("password_12345_1122")
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name        string
		email       string
		passwd      string
		credRep     auth.CredentialRepository
		expected    auth.Credential
		expectedErr error
	}{
		{
			name:   "success",
			email:  "example@example.org",
			passwd: "password_12345_1122",
			credRep: &mock.CredentialRepository{
				ByEmailFn: func(ctx context.Context, email string) (auth.Credential, error) {
					return auth.Credential{
						ID:       1,
						Password: hash,
						Token:    token,
						Email:    "example@example.org",
					}, nil
				},
			},
			expected: auth.Credential{
				ID:       1,
				Password: hash,
				Token:    token,
				Email:    "example@example.org",
			},
		},
		{
			name:   "error - user not found",
			email:  "example@example.org",
			passwd: "password_12345_1122",
			credRep: &mock.CredentialRepository{
				ByEmailFn: func(ctx context.Context, email string) (auth.Credential, error) {
					return auth.Credential{}, auth.NewError(auth.ErrCredNotFound, "Credential not found")
				},
			},
			expectedErr: auth.WrapError(
				auth.NewError(auth.ErrCredNotFound, "Credential not found"),
				auth.ErrAuth,
				"Auth failed",
			),
		},
		{
			name:   "error - different passwords",
			email:  "example@example.org",
			passwd: "12345",
			credRep: &mock.CredentialRepository{
				ByEmailFn: func(ctx context.Context, email string) (auth.Credential, error) {
					return auth.Credential{
						ID:       1,
						Password: hash,
						Token:    token,
						Email:    "example@example.org",
					}, nil
				},
			},
			expectedErr: auth.NewError(auth.ErrAuth, "Auth failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewCredentialService(
				tc.credRep,
				nowFunc,
				func(n int) (string, error) {
					return token, nil
				},
			)
			cred, err := s.Auth(context.Background(), tc.email, tc.passwd)
			if err != nil && tc.expectedErr == nil {
				t.Fatal(err)
			}
			if err != nil {
				if err != tc.expectedErr {
					t.Fatalf("bad error, expected: %v, got: %v", tc.expectedErr, err)
				}
			}

			if diff := cmp.Diff(tc.expected, cred); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
