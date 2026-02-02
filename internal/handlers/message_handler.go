package handlers

import (
	"chats/internal/domain"
	"chats/internal/helpers"
	"chats/internal/services"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

type MessageHandler struct {
	service services.MessageService
}

func NewMessageHandler(service services.MessageService) *MessageHandler {
	return &MessageHandler{
		service: service,
	}
}

func (h *MessageHandler) HandleCreateMessage(w http.ResponseWriter, r *http.Request) {
	logger := slog.Default()

	if r.Method != http.MethodPost {
		logger.Warn("method not allowed", "method", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	id, err := helpers.ExtractIDFromPath(r)
	if err != nil {
		logger.Warn("Bad Request", "error", err)
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Warn("Bad Request", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	request.Text = strings.TrimSpace(request.Text)
	if request.Text == "" {
		logger.Warn("Bad Request", "error", "Text cannot be empty")
		http.Error(w, "Text cannot be empty", http.StatusBadRequest)
		return
	}

	if len(request.Text) > 5000 {
		logger.Warn("Bad Request", "error", "Text too long")
		http.Error(w, "Text must be 5000 characters or less", http.StatusBadRequest)
		return
	}

	message, err := h.service.CreateMessage(r.Context(), id, request.Text)
	if err != nil {
		logger.Error("Error creating message", "error", err)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			http.Error(w, "Chat not found", http.StatusNotFound)
		case errors.Is(err, domain.ErrInvalidInput):
			http.Error(w, "Invalid input", http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		logger.Error("Error encoding response", "error", err)
	}
}
