package app

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/Krokozabra213/test_api/internal/domain"
	"github.com/Krokozabra213/test_api/tests/app/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func TestCreateChat(t *testing.T) {
	ctx, st := suite.New(t)
	t.Cleanup(func() {
		st.CleanupTestData()
	})

	title := "Test Chat"
	resp, err := st.HTTPClient.POST(ctx, "/chats", map[string]string{
		"title": title,
	})
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var chat domain.Chat
	err = resp.JSON(&chat)
	require.NoError(t, err)

	assert.Equal(t, title, chat.Title)
	assert.NotZero(t, chat.ID)
	assert.NotZero(t, chat.CreatedAt)
}

func TestCreateChat_Validate(t *testing.T) {
	ctx, st := suite.New(t)
	t.Cleanup(func() {
		st.CleanupTestData()
	})

	title := strings.Repeat("t", 220)
	resp, err := st.HTTPClient.POST(ctx, "/chats", map[string]string{
		"title": title,
	})
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var errorResp ErrorResponse
	err = resp.JSON(&errorResp)
	require.NoError(t, err)

	assert.Contains(t, errorResp.Error, "title should be between")
}

func TestCreateChat_Sanitize(t *testing.T) {
	ctx, st := suite.New(t)
	t.Cleanup(func() {
		st.CleanupTestData()
	})

	title := strings.Repeat(" ", 200) + strings.Repeat("t", 10) + strings.Repeat(" ", 200)
	resp, err := st.HTTPClient.POST(ctx, "/chats", map[string]string{
		"title": title,
	})
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var chat domain.Chat
	err = resp.JSON(&chat)
	require.NoError(t, err)

	assert.Equal(t, strings.Repeat("t", 10), chat.Title)
	assert.NotZero(t, chat.ID)
	assert.NotZero(t, chat.CreatedAt)
}

func TestSendMessage(t *testing.T) {
	ctx, st := suite.New(t)
	t.Cleanup(func() {
		st.CleanupTestData()
	})

	title := "Test Chat"
	resp, err := st.HTTPClient.POST(ctx, "/chats", map[string]string{
		"title": title,
	})
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var chat domain.Chat
	err = resp.JSON(&chat)
	require.NoError(t, err)

	text := "Test Text"
	path := fmt.Sprintf("/chats/%d/messages", chat.ID)

	resp, err = st.HTTPClient.POST(ctx, path, map[string]string{
		"text": text,
	})
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var message domain.Message
	err = resp.JSON(&message)
	require.NoError(t, err)

	assert.Equal(t, chat.ID, message.ChatID)
	assert.Equal(t, text, message.Text)
	assert.NotZero(t, message.ID)
	assert.NotZero(t, message.CreatedAt)
}
