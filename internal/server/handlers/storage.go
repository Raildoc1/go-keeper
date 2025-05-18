package handlers

import (
	"context"
	"go-keeper/internal/common/protocol"
	"go-keeper/internal/server/dto"
	"go-keeper/pkg/logging"
	"net/http"

	"go.uber.org/zap"
)

type StorageService interface {
	Store(ctx context.Context, entry dto.Entry) error
	Load(ctx context.Context, id string) (dto.Entry, error)
	Delete(ctx context.Context, id string) error
	LoadAll(ctx context.Context) ([]dto.Entry, error)
}

type StorageHandler struct {
	service StorageService
	logger  *logging.ZapLogger
}

func newStorageHandler(service StorageService, logger *logging.ZapLogger) *StorageHandler {
	return &StorageHandler{
		service: service,
		logger:  logger,
	}
}

// Store Handler.

type StoreHandler struct {
	*StorageHandler
}

func NewStoreHandler(service StorageService, logger *logging.ZapLogger) *StoreHandler {
	return &StoreHandler{
		StorageHandler: newStorageHandler(service, logger),
	}
}

func (h *StoreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer closeBody(r.Context(), r.Body, h.logger)

	input, err := decodeJSON[protocol.Entry](r.Body)
	if err != nil {
		h.logger.DebugCtx(r.Context(), "input decoding error", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	entry := dto.Entry{
		ID:       input.ID,
		Metadata: input.Metadata,
		Data:     input.Data,
	}
	err = h.service.Store(r.Context(), entry)
	if err != nil {
		switch {
		default:
			h.logger.ErrorCtx(r.Context(), "storage service error", zap.Error(err), zap.Any("input", input))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// Load Handler.

type LoadHandler struct {
	*StorageHandler
}

func NewLoadHandler(service StorageService, logger *logging.ZapLogger) *LoadHandler {
	return &LoadHandler{
		StorageHandler: newStorageHandler(service, logger),
	}
}

func (h *LoadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer closeBody(r.Context(), r.Body, h.logger)

	input, err := decodeJSON[protocol.LoadRequest](r.Body)
	if err != nil {
		h.logger.DebugCtx(r.Context(), "input decoding error", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	entry, err := h.service.Load(r.Context(), input.ID)
	if err != nil {
		switch {
		default:
			h.logger.ErrorCtx(r.Context(), "storage service error", zap.Error(err), zap.Any("input", input))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err = encodeJSON(w, entry)
	if err != nil {
		h.logger.ErrorCtx(r.Context(), "error writing result", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Load All Handler.

type LoadAllHandler struct {
	*StorageHandler
}

func NewLoadAllHandler(service StorageService, logger *logging.ZapLogger) *LoadAllHandler {
	return &LoadAllHandler{
		StorageHandler: newStorageHandler(service, logger),
	}
}

func (h *LoadAllHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer closeBody(r.Context(), r.Body, h.logger)

	entries, err := h.service.LoadAll(r.Context())
	if err != nil {
		switch {
		default:
			h.logger.ErrorCtx(r.Context(), "storage service error", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err = encodeJSON(w, entries)
	if err != nil {
		h.logger.ErrorCtx(r.Context(), "error writing result", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
