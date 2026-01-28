package handlers

import (
	"chats/internal/helpers"
	"chats/internal/services"
	"encoding/json"
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
	}

	id, err := helpers.ExtractIDFromPath(r)
	if err != nil {
		logger.Warn("Bad Request", "error", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var request struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Warn("Bad Request", "error", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	request.Text = strings.TrimSpace(request.Text)
	if request.Text == "" || len(request.Text) > 5000 {
		logger.Warn("Bad Request", "error", "Incorrect text length")
		http.Error(w, "Incorrect text length", http.StatusBadRequest)
		return
	}

	message, err := h.service.CreateMessage(r.Context(), id, request.Text)
	if err != nil {
		logger.Warn("Bad Request", "error", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		logger.Warn("Bad Request", "error", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}
