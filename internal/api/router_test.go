package api

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	auth "github.com/kl09/auth-go"
	"github.com/kl09/auth-go/internal/mock"
)

var (
	now     = time.Date(2020, time.April, 15, 10, 11, 12, 0, time.UTC)
	nowFunc = func() time.Time {
		return now
	}
)

func TestUser_ByToken(t *testing.T) {
	cases := []struct {
		name       string
		token      string
		wantResp   string
		wantStatus int
		credRep    auth.CredentialRepository
	}{
		{
			name:       "success",
			token:      "12345",
			wantResp:   `{"id":1,"token":"token","email":"example@example.org","email_tmp":"","email_verified":false,"created_at":"2020-04-15T10:11:12Z","updated_at":"2020-04-15T10:11:12Z"}` + "\n",
			wantStatus: http.StatusOK,
			credRep: &mock.CredentialRepositoryMock{
				ByTokenFunc: func(ctx context.Context, token string) (auth.Credential, error) {
					return auth.Credential{
						ID:        1,
						Password:  "12345",
						Email:     "example@example.org",
						Token:     "token",
						CreatedAt: now,
						UpdatedAt: now,
					}, nil
				},
			},
		},
		{
			name:       "error - token not found",
			token:      "12345",
			wantResp:   `{"error":{"code":"credential_not_found","message":"Credential not found"}}` + "\n",
			wantStatus: http.StatusNotFound,
			credRep: &mock.CredentialRepositoryMock{
				ByTokenFunc: func(ctx context.Context, token string) (auth.Credential, error) {
					return auth.Credential{}, auth.NewError(auth.ErrCredNotFound, "Credential not found")
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewRouter(NewCredentialService(tc.credRep, nowFunc, nil)).Handler().Server.Handler

			srv := httptest.NewServer(h)
			defer srv.Close()

			req, err := http.NewRequest(
				"GET",
				fmt.Sprintf("%s/v1/users-by-token/%s", srv.URL, tc.token),
				strings.NewReader(``),
			)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if diff := cmp.Diff(tc.wantStatus, resp.StatusCode); diff != "" {
				t.Error(diff)
			}

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.wantResp, string(b)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestUser_Register(t *testing.T) {
	token := "1234abcd"

	cases := []struct {
		name        string
		requestBody string
		wantResp    string
		wantStatus  int
		credRep     auth.CredentialRepository
	}{
		{
			name:        "success",
			requestBody: `{"email":"example@example.org","password":"66554433"}`,
			wantResp:    `{"id":1,"token":"1234abcd","email":"example@example.org","email_tmp":"","email_verified":false,"created_at":"2020-04-15T10:11:12Z","updated_at":"2020-04-15T10:11:12Z"}` + "\n",
			wantStatus:  http.StatusOK,
			credRep: &mock.CredentialRepositoryMock{
				CreateFunc: func(ctx context.Context, c *auth.Credential) error {
					c.ID = 1
					return nil
				},
				ByEmailFunc: func(ctx context.Context, email string) (auth.Credential, error) {
					return auth.Credential{}, auth.NewError(auth.ErrCredNotFound, "Credential not found")
				},
			},
		},
		{
			name:        "error - internal error",
			requestBody: `{"email":"example@example.org","password":"66554433"}`,
			wantResp:    `{"error":{"code":"internal","message":"no message"}}` + "\n",
			wantStatus:  http.StatusInternalServerError,
			credRep: &mock.CredentialRepositoryMock{
				CreateFunc: func(ctx context.Context, c *auth.Credential) error {
					return errors.New("some error")
				},
				ByEmailFunc: func(ctx context.Context, email string) (auth.Credential, error) {
					return auth.Credential{}, auth.NewError(auth.ErrCredNotFound, "Credential not found")
				},
			},
		},
		{
			name:        "user already exists",
			requestBody: `{"email":"example@example.org","password":"66554433"}`,
			wantResp:    `{"error":{"code":"email_already_exists","message":"User with this email already exists."}}` + "\n",
			wantStatus:  http.StatusInternalServerError,
			credRep: &mock.CredentialRepositoryMock{
				CreateFunc: func(ctx context.Context, c *auth.Credential) error {
					t.Fatal("method shouldn't be called")
					return nil
				},
				ByEmailFunc: func(ctx context.Context, email string) (auth.Credential, error) {
					return auth.Credential{}, nil
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewRouter(NewCredentialService(
				tc.credRep,
				nowFunc,
				func(n int) (string, error) {
					return token, nil
				},
			)).Handler().Server.Handler

			srv := httptest.NewServer(h)
			defer srv.Close()

			req, err := http.NewRequest(
				"POST",
				fmt.Sprintf("%s/v1/register", srv.URL),
				strings.NewReader(tc.requestBody),
			)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if diff := cmp.Diff(tc.wantStatus, resp.StatusCode); diff != "" {
				t.Error(diff)
			}

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.wantResp, string(b)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestUser_Auth(t *testing.T) {
	token := "1234abcd"

	hash, err := hashAndSalt("66554433")
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name        string
		requestBody string
		wantResp    string
		wantStatus  int
		credRep     auth.CredentialRepository
	}{
		{
			name:        "success",
			requestBody: `{"email":"example@example.org","password":"66554433"}`,
			wantResp:    `{"id":1,"token":"token","email":"example@example.org","email_tmp":"","email_verified":false,"created_at":"2020-04-15T10:11:12Z","updated_at":"2020-04-15T10:11:12Z"}` + "\n",
			wantStatus:  http.StatusOK,
			credRep: &mock.CredentialRepositoryMock{
				ByEmailFunc: func(ctx context.Context, email string) (auth.Credential, error) {
					return auth.Credential{
						ID:        1,
						Password:  hash,
						Email:     "example@example.org",
						Token:     "token",
						CreatedAt: now,
						UpdatedAt: now,
					}, nil
				},
			},
		},
		{
			name:        "error - bad password",
			requestBody: `{"email":"example@example.org","password":"66554433"}`,
			wantResp:    `{"error":{"code":"auth_failed","message":"Auth failed"}}` + "\n",
			wantStatus:  http.StatusUnauthorized,
			credRep: &mock.CredentialRepositoryMock{
				ByEmailFunc: func(ctx context.Context, email string) (auth.Credential, error) {
					return auth.Credential{
						ID:        1,
						Password:  "66554433",
						Email:     "example@example.org",
						Token:     "token",
						CreatedAt: now,
						UpdatedAt: now,
					}, nil
				},
			},
		},
		{
			name:        "error - not found by email",
			requestBody: `{"email":"example@example.org","password":"66554433"}`,
			wantResp:    `{"error":{"code":"auth_failed","message":"Auth failed"}}` + "\n",
			wantStatus:  http.StatusUnauthorized,
			credRep: &mock.CredentialRepositoryMock{
				ByEmailFunc: func(ctx context.Context, email string) (auth.Credential, error) {
					return auth.Credential{}, auth.NewError(auth.ErrCredNotFound, "Credential not found")
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewRouter(NewCredentialService(
				tc.credRep,
				nowFunc,
				func(n int) (string, error) {
					return token, nil
				},
			)).Handler().Server.Handler

			srv := httptest.NewServer(h)
			defer srv.Close()

			req, err := http.NewRequest(
				"POST",
				fmt.Sprintf("%s/v1/auth", srv.URL),
				strings.NewReader(tc.requestBody),
			)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if diff := cmp.Diff(tc.wantStatus, resp.StatusCode); diff != "" {
				t.Error(diff)
			}

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.wantResp, string(b)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func Test_404_error(t *testing.T) {
	h := NewRouter(NewCredentialService(
		nil,
		nowFunc,
		func(n int) (string, error) {
			return "", nil
		},
	)).Handler().Server.Handler

	srv := httptest.NewServer(h)
	defer srv.Close()

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/bad_url", srv.URL),
		strings.NewReader(``),
	)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatal("bad status code")
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"error":{"code":"http_404","message":"Not Found"}}` + "\n"
	if diff := cmp.Diff(string(b), expected); diff != "" {
		t.Fatal(diff)
	}
}
