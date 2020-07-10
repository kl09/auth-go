package pg_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	auth "github.com/kl09/auth-go"
	"github.com/kl09/auth-go/internal/pg"
)

func TestCredentialRepository_Create(t *testing.T) {
	c := setUp(t)
	defer c.Close()

	r := pg.NewCredentialRepository(c)

	now := time.Now()
	cred := auth.Credential{
		Password:  "12345",
		Email:     "example@example.org",
		CreatedAt: now,
		UpdatedAt: now,
	}
	require.Nil(t, r.Create(context.Background(), &cred))

	assert.Equal(t,
		auth.Credential{
			ID:                       1,
			Password:                 "12345",
			Email:                    "example@example.org",
			EmailTmp:                 "",
			EmailVerified:            false,
			VerificationCode:         "",
			VerificationCodeAttempts: 0,
			CreatedAt:                now,
			UpdatedAt:                now,
		},
		cred,
	)
}

func TestCredentialRepository_ByToken(t *testing.T) {
	c := setUp(t)
	defer c.Close()

	r := pg.NewCredentialRepository(c)

	now := time.Date(2020, time.April, 15, 0, 0, 0, 0, time.UTC)
	cred := auth.Credential{
		Password:  "12345",
		Email:     "example@example.org",
		Token:     "token",
		CreatedAt: now,
		UpdatedAt: now,
	}
	require.Nil(t, r.Create(context.Background(), &cred))

	testCases := []struct {
		name          string
		token         string
		expectedCred  auth.Credential
		expectedError error
	}{
		{
			name:  "success get",
			token: "token",
			expectedCred: auth.Credential{
				ID:        1,
				Password:  "12345",
				Email:     "example@example.org",
				Token:     "token",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name:          "error - not found",
			token:         "bad_token",
			expectedError: auth.NewError(auth.ErrCredNotFound, "Credential not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cred, err := r.ByToken(context.Background(), tc.token)

			if diff := cmp.Diff(tc.expectedCred, cred); diff != "" {
				t.Fatal(diff)
			}
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestCredentialRepository_ByID(t *testing.T) {
	c := setUp(t)
	defer c.Close()

	r := pg.NewCredentialRepository(c)

	now := time.Date(2020, time.April, 15, 0, 0, 0, 0, time.UTC)
	cred := auth.Credential{
		Password:  "12345",
		Email:     "example@example.org",
		CreatedAt: now,
		UpdatedAt: now,
	}
	require.Nil(t, r.Create(context.Background(), &cred))

	testCases := []struct {
		name          string
		id            int
		expectedCred  auth.Credential
		expectedError error
	}{
		{
			name: "success get",
			id:   1,
			expectedCred: auth.Credential{
				ID:        1,
				Password:  "12345",
				Email:     "example@example.org",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name:          "error - not found",
			id:            2,
			expectedError: auth.NewError(auth.ErrCredNotFound, "Credential not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cred, err := r.ByID(context.Background(), tc.id)

			if diff := cmp.Diff(tc.expectedCred, cred); diff != "" {
				t.Fatal(diff)
			}
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestCredentialRepository_ByEmail(t *testing.T) {
	c := setUp(t)
	defer c.Close()

	r := pg.NewCredentialRepository(c)

	now := time.Date(2020, time.April, 15, 0, 0, 0, 0, time.UTC)
	cred := auth.Credential{
		Password:  "12345",
		Email:     "example@example.org",
		CreatedAt: now,
		UpdatedAt: now,
	}
	require.Nil(t, r.Create(context.Background(), &cred))

	testCases := []struct {
		name          string
		email         string
		expectedCred  auth.Credential
		expectedError error
	}{
		{
			name:  "success get",
			email: "example@example.org",
			expectedCred: auth.Credential{
				ID:        1,
				Password:  "12345",
				Email:     "example@example.org",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name:          "error - not found",
			email:         "example2@example.org",
			expectedError: auth.NewError(auth.ErrCredNotFound, "Credential not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cred, err := r.ByEmail(context.Background(), tc.email)

			if diff := cmp.Diff(tc.expectedCred, cred); diff != "" {
				t.Fatal(diff)
			}
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
