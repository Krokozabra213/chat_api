package domain

import "time"

// DTO for output response (to the client)
type ChatMessageOutput struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Messages  []Message `json:"messages"`
}

func NewChatMessageOutput(chatID int64, title string, createdAt time.Time, messages []Message) *ChatMessageOutput {
	return &ChatMessageOutput{
		ID:        chatID,
		Title:     title,
		CreatedAt: createdAt,
		Messages:  messages,
	}
}
