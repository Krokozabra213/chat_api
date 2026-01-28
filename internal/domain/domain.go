package domain

import "time"

// Chat represents a chat entity used both as a DTO for client communication
// and as a business logic model.
type Chat struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

func NewChat(title string) Chat {
	return Chat{
		Title: title,
	}
}

// Message represents a message entity used both as a DTO for client communication
// and as a business logic model.
type Message struct {
	ID        int64     `json:"id"`
	ChatID    int64     `json:"chat_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func NewMessage(chatID int64, text string) Message {
	return Message{
		ChatID: chatID,
		Text:   text,
	}
}
