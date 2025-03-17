package main

import (
	"chat/cmd/database"
	router "chat/cmd/routes"

	"chat/internal/handler"
	"chat/internal/repository"
	"chat/internal/usecase"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠ No .env file found, using system environment variables")
	}

	database.ConnectDB()

	//auth
	userRepo := repository.NewUserRepository(database.DB)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	//roomchat
	roomChatRepo := repository.NewRoomChatRepository(database.DB)
	roomChatUseCase := usecase.NewRoomChatUseCase(roomChatRepo)
	roomChatHandler := handler.NewRoomChatUserHandler(roomChatUseCase)

	//chat
	chatRepo := repository.NewChatRepository(database.DB)
	chatUseCase := usecase.NewChatUsecase(chatRepo, roomChatRepo)
	chatHandler := handler.NewChatHandler(chatUseCase)

	// Setup Router
	r := router.SetupRoutes(userHandler, roomChatHandler, chatHandler)

	// Mulai Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("🚀 Server running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
