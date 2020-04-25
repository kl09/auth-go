package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	auth "github.com/kl09/auth-go"
)

type response struct {
	ID            int       `json:"id"`
	Token         string    `json:"token"`
	Email         string    `json:"email"`
	EmailTmp      string    `json:"email_tmp"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func credToResponse(cred auth.Credential) response {
	return response{
		ID:            cred.ID,
		Token:         cred.Token,
		Email:         cred.Email,
		EmailTmp:      cred.EmailTmp,
		EmailVerified: cred.EmailVerified,
		CreatedAt:     cred.CreatedAt,
		UpdatedAt:     cred.UpdatedAt,
	}
}

// userByToken retrieves the user by token.
func (r *Router) userByToken(c echo.Context) error {
	token := c.Param("token")

	cred, err := r.credService.ByToken(c.Request().Context(), token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, credToResponse(cred))
}

// registerUser creates a new user.
func (r *Router) registerUser(c echo.Context) error {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := c.Bind(&request)
	if err != nil {
		return err
	}

	cred := auth.Credential{
		Email:    request.Email,
		Password: request.Password,
	}

	err = r.credService.Register(c.Request().Context(), &cred)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, credToResponse(cred))
}

func (r *Router) auth(c echo.Context) error {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := c.Bind(&request)
	if err != nil {
		return err
	}

	cred, err := r.credService.Auth(c.Request().Context(), request.Email, request.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, credToResponse(cred))
}

func customHTTPErrorHandler(err error, c echo.Context) {
	httpStatus := http.StatusInternalServerError
	errResp := struct {
		Err *auth.Error `json:"error"`
	}{&auth.Error{
		Code:    auth.ErrorCode(err),
		Message: auth.ErrorMsg(err),
	}}

	switch errI := err.(type) {
	case *echo.HTTPError:
		httpStatus = errI.Code
		errResp.Err.Code = fmt.Sprintf("http_%d", errI.Code)

		if msg, ok := errI.Message.(string); ok {
			errResp.Err.Message = msg
		}
	case auth.Error:
		switch errI.Code {
		case auth.ErrCredNotFound:
			httpStatus = http.StatusNotFound
		case auth.ErrAuth:
			httpStatus = http.StatusUnauthorized
		}
	default:
		c.Logger().Error(err)
	}

	err = c.JSON(httpStatus, errResp)
	if err != nil {
		c.Logger().Error(err)
	}
}
