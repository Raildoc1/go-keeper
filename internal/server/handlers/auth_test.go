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

	loginHandlerSetup := handlerSetup{
		handler: NewLoginHandler(authService, logger),
		method:  http.MethodPost,
		url:     "/api/user/login",
	}

	registerHandlerSetup := handlerSetup{
		handler: NewRegisterHandler(authService, logger),
		method:  http.MethodPost,
		url:     "/api/user/register",
	}

	tests := []handlerTestData{
		{
			testName:       "Register Success",
			handlerSetup:   registerHandlerSetup,
			body:           testutils.MustCreateCredsJSON("test_user", "test_password"),
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Login Success",
			handlerSetup:   loginHandlerSetup,
			body:           testutils.MustCreateCredsJSON("test_user", "test_password"),
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Register Fail Same User",
			handlerSetup:   registerHandlerSetup,
			body:           testutils.MustCreateCredsJSON("test_user", "test_password"),
			expectedStatus: http.StatusConflict,
		},
		{
			testName:       "Login Fail Non-Existent User",
			handlerSetup:   loginHandlerSetup,
			body:           testutils.MustCreateCredsJSON("wrong_user", "test_password"),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			testName:       "Login Fail Wrong Password",
			handlerSetup:   loginHandlerSetup,
			body:           testutils.MustCreateCredsJSON("test_user", "wrong_password"),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			testName:       "Login Fail Invalid Input",
			handlerSetup:   loginHandlerSetup,
			body:           "} invalid json {",
			expectedStatus: http.StatusBadRequest,
		},
		{
			testName:       "Register Fail Invalid Input",
			handlerSetup:   registerHandlerSetup,
			body:           "} invalid json {",
			expectedStatus: http.StatusBadRequest,
		},
	}

	performHTTPHandlerTests(t, tests)
}
