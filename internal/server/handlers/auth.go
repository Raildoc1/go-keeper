package handlers

import (
	"context"
	"errors"
	"go-keeper/internal/common/protocol"
	"go-keeper/internal/server/dto"
	"go-keeper/internal/server/services"
	"go-keeper/pkg/logging"
	"net/http"

	"go.uber.org/zap"
)

type AuthService interface {
	Login(ctx context.Context, creds dto.Creds) (string, error)
	Register(ctx context.Context, creds dto.Creds) (string, error)
}

// Username Handler.

type LoginHandler struct {
	service AuthService
	logger  *logging.ZapLogger
}

func NewLoginHandler(service AuthService, logger *logging.ZapLogger) *LoginHandler {
	return &LoginHandler{
		service: service,
		logger:  logger,
	}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer closeBody(r.Context(), r.Body, h.logger)

	input, err := decodeJSON[protocol.Creds](r.Body)
	if err != nil {
		h.logger.DebugCtx(r.Context(), "input decoding error", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	creds := dto.Creds{
		Username: input.Username,
		Password: input.Password,
	}
	tkn, err := h.service.Login(r.Context(), creds)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			h.logger.DebugCtx(r.Context(), err.Error(), zap.Any("input", input))
			w.WriteHeader(http.StatusUnauthorized)
			return
		default:
			h.logger.ErrorCtx(r.Context(), "login service error", zap.Error(err), zap.Any("input", input))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Authorization", "Bearer "+tkn)
}

// Registration Handler.

type RegisterHandler struct {
	service AuthService
	logger  *logging.ZapLogger
}

func NewRegisterHandler(service AuthService, logger *logging.ZapLogger) *RegisterHandler {
	return &RegisterHandler{
		service: service,
		logger:  logger,
	}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer closeBody(r.Context(), r.Body, h.logger)

	input, err := decodeJSON[protocol.Creds](r.Body)
	if err != nil {
		h.logger.DebugCtx(r.Context(), "input decoding error", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	creds := dto.Creds{
		Username: input.Username,
		Password: input.Password,
	}
	tkn, err := h.service.Register(r.Context(), creds)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrLoginTaken):
			h.logger.DebugCtx(r.Context(), err.Error(), zap.String("login", input.Username))
			w.WriteHeader(http.StatusConflict)
			return
		default:
			h.logger.ErrorCtx(r.Context(), "registration service error", zap.Error(err), zap.Any("input", input))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Authorization", "Bearer "+tkn)
}
