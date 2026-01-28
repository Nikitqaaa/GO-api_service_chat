package route

import (
	"chats/internal/handlers"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupQuestionRoutes(chatHandler *handlers.ChatHandler, messageHandler *handlers.MessageHandler) http.Handler {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Get("/chats/{id}", chatHandler.HandleGetChat)
		r.Post("/chats", chatHandler.HandleCreateChat)
		r.Delete("/chats/{id}", chatHandler.HandleDeleteChat)

		r.Post("/chats/{id}/messages", messageHandler.HandleCreateMessage)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Health check from %s", r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		if err != nil {
			log.Printf("Error encoding response: %v", err)
			return
		}
	})

	return r
}
