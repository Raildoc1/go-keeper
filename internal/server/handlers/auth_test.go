package handlers

import (
	"go-keeper/internal/server/services"
	"go-keeper/internal/server/testutils"
	"go-keeper/internal/server/testutils/mock/repositories"
	"go-keeper/pkg/logging"
	"net/http"
	"testing"
)

func TestAuthHandlers(t *testing.T) {
	authRepository := repositories.NewAuthRepositoryMock()
	tokenFactory := repositories.NewTokenFactoryMock()

	authService := services.NewAuthService(authRepository, tokenFactory)

	logger := logging.NewNopLogger()

	loginHandlerSetup := testutils.HandlerSetup{
		Handler: NewLoginHandler(authService, logger),
		Method:  http.MethodPost,
		URL:     "/api/user/login",
	}

	registerHandlerSetup := testutils.HandlerSetup{
		Handler: NewRegisterHandler(authService, logger),
		Method:  http.MethodPost,
		URL:     "/api/user/register",
	}

	tests := []testutils.HandlerTestData{
		{
			TestName:       "Register Success",
			HandlerSetup:   registerHandlerSetup,
			Body:           testutils.MustCreateCredsJSON("test_user", "test_password"),
			ExpectedStatus: http.StatusOK,
		},
		{
			TestName:       "Login Success",
			HandlerSetup:   loginHandlerSetup,
			Body:           testutils.MustCreateCredsJSON("test_user", "test_password"),
			ExpectedStatus: http.StatusOK,
		},
		{
			TestName:       "Register Fail Same User",
			HandlerSetup:   registerHandlerSetup,
			Body:           testutils.MustCreateCredsJSON("test_user", "test_password"),
			ExpectedStatus: http.StatusConflict,
		},
		{
			TestName:       "Login Fail Non-Existent User",
			HandlerSetup:   loginHandlerSetup,
			Body:           testutils.MustCreateCredsJSON("wrong_user", "test_password"),
			ExpectedStatus: http.StatusUnauthorized,
		},
		{
			TestName:       "Login Fail Wrong Password",
			HandlerSetup:   loginHandlerSetup,
			Body:           testutils.MustCreateCredsJSON("test_user", "wrong_password"),
			ExpectedStatus: http.StatusUnauthorized,
		},
		{
			TestName:       "Login Fail Invalid Input",
			HandlerSetup:   loginHandlerSetup,
			Body:           "} invalid json {",
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			TestName:       "Register Fail Invalid Input",
			HandlerSetup:   registerHandlerSetup,
			Body:           "} invalid json {",
			ExpectedStatus: http.StatusBadRequest,
		},
	}

	testutils.PerformHTTPHandlerTests(t, tests)
}
