package handlers

import (
	"chats/internal/domain"
	"chats/internal/helpers"
	"chats/internal/services"
	"encoding/json"
	"log/slog"
	"net/http"
)

type ChatHandler struct {
	service services.ChatService
}

func NewChatHandler(service services.ChatService) *ChatHandler {
	return &ChatHandler{
		service: service,
	}
}

func (h *ChatHandler) HandleCreateChat(w http.ResponseWriter, r *http.Request) {
	logger := slog.Default()

	if r.Method != http.MethodPost {
		logger.Warn("method not allowed", "method", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var req domain.Chat

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Bad Request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	createdChat, err := h.service.CreateChat(r.Context(), req.Title)
	if err != nil {
		logger.Warn("Internal Server Error", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdChat)
	if err != nil {
		logger.Error("Error encoding response:", "error", err)
		return
	}
}

func (h *ChatHandler) HandleGetChat(w http.ResponseWriter, r *http.Request) {
	logger := slog.Default()

	if r.Method != http.MethodGet {
		logger.Warn("method not allowed", "method", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	id, err := helpers.ExtractIDFromPath(r)
	if err != nil {
		logger.Warn("Bad Request", "error", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	limit := helpers.ParseLimitParam(r, 20, 100)
	chat, err := h.service.GetChat(r.Context(), id, limit)
	if err != nil {
		logger.Warn("Not Found", "error", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(chat)
	if err != nil {
		logger.Error("Error encoding response:", "error", err)
		return
	}
}

func (h *ChatHandler) HandleDeleteChat(w http.ResponseWriter, r *http.Request) {
	logger := slog.Default()
	if r.Method != http.MethodDelete {
		logger.Warn("method not allowed", "method", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	id, err := helpers.ExtractIDFromPath(r)
	if err != nil {
		logger.Warn("Bad Request", "error", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteChat(r.Context(), id)
	if err != nil {
		logger.Warn("Not Found", "error", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
