// Package handler provides HTTP handlers for API.
package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Krokozabra213/test_api/internal/domain"
	"github.com/Krokozabra213/test_api/pkg/request"
)

// Business defines business layer interface.
type Business interface {
	CreateChat(ctx context.Context, title string) (*domain.Chat, error)
	DeleteChat(ctx context.Context, chatID int64) error
	CreateMessage(ctx context.Context, chatID int64, text string) (*domain.Message, error)
	ReadChatMessages(ctx context.Context, chatID int64, limit int) (*domain.ChatMessageOutput, error)
}

// Handler handles HTTP requests.
type Handler struct {
	log      *slog.Logger
	business Business
}

// NewHandler creates a new Handler and registers routes.
func New(router *http.ServeMux, log *slog.Logger, business Business) {
	handler := &Handler{
		log:      log,
		business: business,
	}
	router.HandleFunc("POST /chats", handler.CreateChat())
	router.HandleFunc("POST /chats/{id}/messages", handler.SendMessage())
	router.HandleFunc("GET /chats/{id}", handler.GetChatMessages())
	router.HandleFunc("DELETE /chats/{id}", handler.DeleteChat())
}

// CreateChat handles chat creation.
func (h *Handler) CreateChat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.DecodeAndValidate[domain.CreateChatInput](r)
		if err != nil {
			h.respondError(w, http.StatusBadRequest, err.Error()) // 400
			return
		}
		body.Sanitize()

		chat, err := h.business.CreateChat(r.Context(), body.Title)
		if err != nil {
			h.handleBusinessError(w, err)
			return
		}

		h.respond(w, http.StatusCreated, chat)
	}
}

// SendMessage handles message creation.
func (h *Handler) SendMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatID, ok := h.parseChatID(w, r)
		if !ok {
			return
		}

		body, err := request.DecodeAndValidate[domain.CreateMessageInput](r)
		if err != nil {
			h.respondError(w, http.StatusBadRequest, err.Error()) // 400
			return
		}
		body.Sanitize()

		message, err := h.business.CreateMessage(r.Context(), chatID, body.Text)
		if err != nil {
			h.handleBusinessError(w, err)
			return
		}

		h.respond(w, http.StatusCreated, message)
	}
}

// GetChatMessages handles getting chat with messages.
func (h *Handler) GetChatMessages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatID, ok := h.parseChatID(w, r)
		if !ok {
			return
		}

		limit := h.parseLimit(r)

		ChatMessage, err := h.business.ReadChatMessages(r.Context(), chatID, limit)
		if err != nil {
			h.handleBusinessError(w, err)
			return
		}

		h.respond(w, http.StatusCreated, ChatMessage)
	}
}

// DeleteChat handles chat deletion.
func (h *Handler) DeleteChat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatID, ok := h.parseChatID(w, r)
		if !ok {
			return
		}

		err := h.business.DeleteChat(r.Context(), chatID)
		if err != nil {
			h.handleBusinessError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
