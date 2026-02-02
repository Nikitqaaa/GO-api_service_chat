package main

import (
	"chats/internal/config"
	"chats/internal/database"
	"chats/internal/handlers"
	"chats/internal/repositories"
	"chats/internal/route"
	"chats/internal/services"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewDatabase(cfg.DB)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	chatRepo := repositories.NewChatRepository(db.DB)
	chatService := services.NewChatService(chatRepo)
	chatHandler := handlers.NewChatHandler(chatService)

	messageRepo := repositories.NewMessageRepository(db.DB)
	messageService := services.NewMessageService(messageRepo, chatService)
	messageHandler := handlers.NewMessageHandler(messageService)

	apiRoute := route.SetupQuestionRoutes(chatHandler, messageHandler)

	serverAddr := cfg.Server.Address + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", serverAddr)

	if err := http.ListenAndServe(serverAddr, apiRoute); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
