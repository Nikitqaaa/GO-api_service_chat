package handlers

import (
	"bytes"
	"chats/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockChatService struct {
	mock.Mock
}

func (m *MockChatService) CreateChat(ctx context.Context, title string) (*domain.Chat, error) {
	args := m.Called(ctx, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Chat), args.Error(1)
}

func (m *MockChatService) GetChat(ctx context.Context, id uint, limit int) (*domain.Chat, error) {
	args := m.Called(ctx, id, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Chat), args.Error(1)
}

func (m *MockChatService) DeleteChat(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockChatService) ValidateChatExists(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestChatHandler_HandleCreateChat_Success(t *testing.T) {
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	expectedChat := &domain.Chat{
		ID:        1,
		Title:     "Новый чат",
		CreatedAt: time.Now(),
	}

	mockService.On("CreateChat", mock.Anything, "Новый чат").Return(expectedChat, nil)

	requestBody, _ := json.Marshal(map[string]string{"title": "Новый чат"})
	req := httptest.NewRequest("POST", "/api/chats", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreateChat(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var response domain.Chat
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedChat.ID, response.ID)
	assert.Equal(t, expectedChat.Title, response.Title)
}

func TestChatHandler_HandleCreateChat_AlreadyExists(t *testing.T) {
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	mockService.On("CreateChat", mock.Anything, "Существующий чат").
		Return(&domain.Chat{}, domain.ErrAlreadyExists)

	requestBody, _ := json.Marshal(map[string]string{"title": "Существующий чат"})
	req := httptest.NewRequest("POST", "/api/chats", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreateChat(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Contains(t, rr.Body.String(), domain.ErrAlreadyExists.Error())
}

func TestChatHandler_HandleCreateChat_InvalidInput(t *testing.T) {
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	mockService.On("CreateChat", mock.Anything, "").Return(&domain.Chat{}, domain.ErrInvalidInput)

	requestBody, _ := json.Marshal(map[string]string{"title": ""})
	req := httptest.NewRequest("POST", "/api/chats", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreateChat(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestChatHandler_HandleCreateChat_MethodNotAllowed(t *testing.T) {
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	req := httptest.NewRequest("GET", "/api/chats", nil)
	rr := httptest.NewRecorder()

	handler.HandleCreateChat(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestChatHandler_HandleGetChat_Success(t *testing.T) {
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	expectedChat := &domain.Chat{
		ID:    1,
		Title: "Тестовый чат",
	}

	mockService.On("GetChat", mock.Anything, uint(1), 20).Return(expectedChat, nil)

	req := httptest.NewRequest("GET", "/api/chats/1", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetChat(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var response domain.Chat
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedChat.ID, response.ID)
	assert.Equal(t, expectedChat.Title, response.Title)
}

func TestChatHandler_HandleGetChat_WithLimit(t *testing.T) {
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	expectedChat := &domain.Chat{ID: 1, Title: "Чат с лимитом"}

	mockService.On("GetChat", mock.Anything, uint(1), 5).Return(expectedChat, nil)

	req := httptest.NewRequest("GET", "/api/chats/1?limit=5", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetChat(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response domain.Chat
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Чат с лимитом", response.Title)
}

func TestChatHandler_HandleGetChat_NotFound(t *testing.T) {
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	mockService.On("GetChat", mock.Anything, uint(999), 20).Return(&domain.Chat{}, domain.ErrNotFound)

	req := httptest.NewRequest("GET", "/api/chats/999", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetChat(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestChatHandler_HandleGetChat_InvalidID(t *testing.T) {
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	req := httptest.NewRequest("GET", "/api/chats/abc", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetChat(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestChatHandler_HandleGetChat_MethodNotAllowed(t *testing.T) {
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	req := httptest.NewRequest("POST", "/api/chats/1", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetChat(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestChatHandler_HandleDeleteChat_Success(t *testing.T) {
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	mockService.On("DeleteChat", mock.Anything, uint(1)).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/chats/1", nil)
	rr := httptest.NewRecorder()

	handler.HandleDeleteChat(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestChatHandler_HandleDeleteChat_NotFound(t *testing.T) {
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	mockService.On("DeleteChat", mock.Anything, uint(999)).Return(domain.ErrNotFound)

	req := httptest.NewRequest("DELETE", "/api/chats/999", nil)
	rr := httptest.NewRecorder()

	handler.HandleDeleteChat(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestChatHandler_HandleDeleteChat_InvalidID(t *testing.T) {
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	req := httptest.NewRequest("DELETE", "/api/chats/abc", nil)
	rr := httptest.NewRecorder()

	handler.HandleDeleteChat(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestChatHandler_HandleDeleteChat_MethodNotAllowed(t *testing.T) {
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	req := httptest.NewRequest("GET", "/api/chats/1", nil)
	rr := httptest.NewRecorder()

	handler.HandleDeleteChat(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestChatHandler_HandleCreateChat_InternalError(t *testing.T) {
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	mockService.On("CreateChat", mock.Anything, "Чат").Return(&domain.Chat{}, errors.New("database error"))

	requestBody, _ := json.Marshal(map[string]string{"title": "Чат"})
	req := httptest.NewRequest("POST", "/api/chats", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreateChat(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestChatHandler_HandleGetChat_InternalError(t *testing.T) {
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	mockService.On("GetChat", mock.Anything, uint(1), 20).Return(&domain.Chat{}, errors.New("database error"))

	req := httptest.NewRequest("GET", "/api/chats/1", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetChat(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestChatHandler_HandleCreateChat_ServiceError(t *testing.T) {
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	mockService.On("CreateChat", mock.Anything, "Ошибка").Return(&domain.Chat{}, errors.New("service error"))

	requestBody, _ := json.Marshal(map[string]string{"title": "Ошибка"})
	req := httptest.NewRequest("POST", "/api/chats", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreateChat(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
