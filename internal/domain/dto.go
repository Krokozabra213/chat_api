// Package domain contains business entities and DTOs.
package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

// Validation limits.
const (
	maxTitleLen       = 200
	maxMessageTextLen = 5000
)

// CreateChatInput represents chat creation request.
type CreateChatInput struct {
	Title string `json:"title"`
}

// Validate checks if chat creation input is valid.
func (i CreateChatInput) Validate() error {
	title := strings.TrimSpace(i.Title)
	titleLen := utf8.RuneCountInString(title)
	if titleLen == 0 || titleLen > maxTitleLen {
		return fmt.Errorf("title should be between 1 and %d characters", maxTitleLen)
	}
	return nil
}

// Sanitize normalizes input data.
func (i *CreateChatInput) Sanitize() {
	i.Title = strings.TrimSpace(i.Title)
}

// DeleteChatInput represents chat deletion request.
type DeleteChatInput struct {
	ID int64 `json:"id"`
}

// Validate checks if chat deletion input is valid.
func (i DeleteChatInput) Validate() error {
	if i.ID <= 0 {
		return errors.New("id should be positive")
	}
	return nil
}

// CreateMessageInput represents message creation request.
type CreateMessageInput struct {
	Text string `json:"text"`
}

// Validate checks if chat deletion input is valid.
func (i CreateMessageInput) Validate() error {
	text := strings.TrimSpace(i.Text)
	textLen := utf8.RuneCountInString(text)
	if textLen == 0 || textLen > maxMessageTextLen {
		return fmt.Errorf("text should be between 1 and %d characters", maxMessageTextLen)
	}
	return nil
}

// Sanitize normalizes input data.
func (i *CreateMessageInput) Sanitize() {
	i.Text = strings.TrimSpace(i.Text)
}

// ChatMessageOutput represents chat with messages response.
type ChatMessageOutput struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Messages  []Message `json:"messages"`
}

// NewChatMessageOutput creates a new ChatMessageOutput instance.
func NewChatMessageOutput(chatID int64, title string, createdAt time.Time, messages []Message) *ChatMessageOutput {
	return &ChatMessageOutput{
		ID:        chatID,
		Title:     title,
		CreatedAt: createdAt,
		Messages:  messages,
	}
}
