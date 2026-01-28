package handler

import (
	"context"
	"net/http"

	"github.com/Krokozabra213/test_api/internal/domain"
)

type Business interface {
	CreateChat(ctx context.Context, title string) (*domain.Chat, error)
	DeleteChat(ctx context.Context, chatID int64) error
	CreateMessage(ctx context.Context, chatID int64, text string) (*domain.Message, error)
	ReadChatMessages(ctx context.Context, chatID int64, limit int) (*domain.ChatMessageOutput, error)
}

type Handler struct {
	business Business
}

func NewLinkHandler(router *http.ServeMux, business Business) {
	handler := &Handler{
		business: business,
	}
	router.HandleFunc("POST /chats/", handler.CreateChat())
	// router.HandleFunc("GET /link/{hash}", handler.GoTo())
	// router.Handle("POST /link", middleware.IsAuthed(handler.Create(), handler.Config))
	// router.HandleFunc("DELETE /link/{id}", handler.Delete())
	// router.Handle("PATCH /link/{id}", middleware.IsAuthed(handler.Update(), handler.Config))
}

func (handler *Handler) CreateChat() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := request.HandleBody[CreateLinkRequest](&w, req)
		if err != nil {
			return
		}

		link := NewLink(body.Url)

		for {
			existLink, _ := handler.LinkRepo.GetByHash(link.Hash)
			if existLink == nil {
				break
			}
			link.generateHash()
		}

		createdLink, err := handler.LinkRepo.Create(link)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp.JsonResp(&w, 201, createdLink)
	}
}

