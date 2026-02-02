package handlers

import (
	"bytes"
	"chats/internal/domain"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) CreateMessage(ctx context.Context, chatID uint, text string) (*domain.Message, error) {
	args := m.Called(ctx, chatID, text)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}

func TestMessageHandler_HandleCreateMessage_Success(t *testing.T) {
	mockService := new(MockMessageService)
	handler := NewMessageHandler(mockService)

	expectedMessage := &domain.Message{
		ID:        1,
		ChatID:    123,
		Text:      "Привет всем!",
		CreatedAt: time.Now(),
	}

	mockService.On("CreateMessage", mock.Anything, uint(123), "Привет всем!").Return(expectedMessage, nil)

	requestBody, _ := json.Marshal(map[string]string{"text": "Привет всем!"})
	req := httptest.NewRequest("POST", "/api/chats/123/messages", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreateMessage(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var response domain.Message
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedMessage.ID, response.ID)
	assert.Equal(t, expectedMessage.ChatID, response.ChatID)
	assert.Equal(t, expectedMessage.Text, response.Text)
}

func TestMessageHandler_HandleCreateMessage_MethodNotAllowed(t *testing.T) {
	mockService := new(MockMessageService)
	handler := NewMessageHandler(mockService)

	req := httptest.NewRequest("GET", "/api/chats/123/messages", nil)
	rr := httptest.NewRecorder()

	handler.HandleCreateMessage(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestMessageHandler_HandleCreateMessage_InvalidID(t *testing.T) {
	mockService := new(MockMessageService)
	handler := NewMessageHandler(mockService)

	req := httptest.NewRequest("POST", "/api/chats/abc/messages", nil)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreateMessage(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestMessageHandler_HandleCreateMessage_InvalidJSON(t *testing.T) {
	mockService := new(MockMessageService)
	handler := NewMessageHandler(mockService)

	req := httptest.NewRequest("POST", "/api/chats/123/messages", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreateMessage(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestMessageHandler_HandleCreateMessage_EmptyText(t *testing.T) {
	mockService := new(MockMessageService)
	handler := NewMessageHandler(mockService)

	requestBody, _ := json.Marshal(map[string]string{"text": ""})
	req := httptest.NewRequest("POST", "/api/chats/123/messages", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreateMessage(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockService.AssertNotCalled(t, "CreateMessage")
}

func TestMessageHandler_HandleCreateMessage_TextTooLong(t *testing.T) {
	mockService := new(MockMessageService)
	handler := NewMessageHandler(mockService)

	longText := string(make([]byte, 5001))
	requestBody, _ := json.Marshal(map[string]string{"text": longText})

	req := httptest.NewRequest("POST", "/api/chats/123/messages", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreateMessage(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockService.AssertNotCalled(t, "CreateMessage")
}

func TestMessageHandler_HandleCreateMessage_TextWithSpaces(t *testing.T) {
	mockService := new(MockMessageService)
	handler := NewMessageHandler(mockService)

	expectedMessage := &domain.Message{
		ID:        1,
		ChatID:    123,
		Text:      "Текст с пробелами",
		CreatedAt: time.Now(),
	}

	mockService.On("CreateMessage", mock.Anything, uint(123), "Текст с пробелами").Return(expectedMessage, nil)

	requestBody, _ := json.Marshal(map[string]string{"text": "  Текст с пробелами  "})
	req := httptest.NewRequest("POST", "/api/chats/123/messages", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreateMessage(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestMessageHandler_HandleCreateMessage_InvalidInput(t *testing.T) {
	mockService := new(MockMessageService)
	handler := NewMessageHandler(mockService)

	mockService.On("CreateMessage", mock.Anything, uint(123), "invalid").
		Return(&domain.Message{}, domain.ErrInvalidInput)

	requestBody, _ := json.Marshal(map[string]string{"text": "invalid"})
	req := httptest.NewRequest("POST", "/api/chats/123/messages", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreateMessage(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestMessageHandler_HandleCreateMessage_EmptyBody(t *testing.T) {
	mockService := new(MockMessageService)
	handler := NewMessageHandler(mockService)

	req := httptest.NewRequest("POST", "/api/chats/123/messages", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreateMessage(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestMessageHandler_HandleCreateMessage_MissingTextField(t *testing.T) {
	mockService := new(MockMessageService)
	handler := NewMessageHandler(mockService)

	requestBody, _ := json.Marshal(map[string]string{"wrong_field": "значение"})
	req := httptest.NewRequest("POST", "/api/chats/123/messages", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreateMessage(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestMessageHandler_HandleCreateMessage_ValidMaxLength(t *testing.T) {
	mockService := new(MockMessageService)
	handler := NewMessageHandler(mockService)

	maxText := string(make([]byte, 5000))
	expectedMessage := &domain.Message{
		ID:        1,
		ChatID:    123,
		Text:      maxText,
		CreatedAt: time.Now(),
	}

	mockService.On("CreateMessage", mock.Anything, uint(123), maxText).Return(expectedMessage, nil)

	requestBody, _ := json.Marshal(map[string]string{"text": maxText})
	req := httptest.NewRequest("POST", "/api/chats/123/messages", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreateMessage(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestMessageHandler_HandleCreateMessage_OnlySpaces(t *testing.T) {
	mockService := new(MockMessageService)
	handler := NewMessageHandler(mockService)

	requestBody, _ := json.Marshal(map[string]string{"text": "     "})
	req := httptest.NewRequest("POST", "/api/chats/123/messages", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreateMessage(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockService.AssertNotCalled(t, "CreateMessage")
}

func TestMessageHandler_HandleCreateMessage_ChatNotFound(t *testing.T) {
	mockService := new(MockMessageService)
	handler := NewMessageHandler(mockService)

	mockService.On("CreateMessage", mock.Anything, uint(999), "Текст").
		Return(&domain.Message{}, domain.ErrNotFound)

	requestBody, _ := json.Marshal(map[string]string{"text": "Текст"})
	req := httptest.NewRequest("POST", "/api/chats/999/messages", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreateMessage(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestMessageHandler_HandleCreateMessage_EncodingError(t *testing.T) {
	mockService := new(MockMessageService)
	handler := NewMessageHandler(mockService)

	expectedMessage := &domain.Message{
		ID:     1,
		ChatID: 123,
		Text:   "Текст",
	}

	mockService.On("CreateMessage", mock.Anything, uint(123), "Текст").Return(expectedMessage, nil)

	requestBody, _ := json.Marshal(map[string]string{"text": "Текст"})
	req := httptest.NewRequest("POST", "/api/chats/123/messages", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreateMessage(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}
