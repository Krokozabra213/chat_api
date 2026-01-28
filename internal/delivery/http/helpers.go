// Package handler provides HTTP handlers for API.
package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Krokozabra213/test_api/internal/business"
)

// Limit constraints
const (
	defaultLimit = 20
	minLimit     = 1
	maxLimit     = 100
)

// respond sends JSON response with logging on error.
func (h *Handler) respond(w http.ResponseWriter, statusCode int, data any) {
	if err := JSONResp(w, statusCode, data); err != nil {
		h.log.Error("failed to send response", "error", err)
	}
}

// respondError sends JSON error response with logging on error.
func (h *Handler) respondError(w http.ResponseWriter, statusCode int, message string) {
	if err := JSONError(w, statusCode, message); err != nil {
		h.log.Error("failed to send error response", "error", err)
	}
}

// handleBusinessError maps business errors to HTTP responses.
func (h *Handler) handleBusinessError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, business.ErrChatNotFound):
		h.respondError(w, http.StatusNotFound, ErrNotFound)
	case errors.Is(err, business.ErrTimeout):
		h.respondError(w, http.StatusGatewayTimeout, ErrRequestTimeout)
	default:
		h.respondError(w, http.StatusInternalServerError, ErrInternal)
	}
}

// parseLimit extracts limit from query params with defaults.
func (h *Handler) parseLimit(r *http.Request) int {
	limitString := r.URL.Query().Get("limit")
	if limitString == "" {
		return defaultLimit
	}

	limit, err := strconv.Atoi(limitString)
	if err != nil {
		h.log.Warn("invalid limit parameter", "value", limitString, "error", err)
		return defaultLimit
	}

	return clamp(limit, minLimit, maxLimit)
}

// parseChatID extracts and validates chat ID from path.
func (h *Handler) parseChatID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	chatIDString := r.PathValue("id")
	chatID, err := strconv.ParseInt(chatIDString, 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, ErrInvalidChatID)
		return 0, false
	}

	if chatID <= 0 {
		h.respondError(w, http.StatusBadRequest, ErrInvalidChatID)
		return 0, false
	}

	return chatID, true
}
