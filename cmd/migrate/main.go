package main

import (
	"chat/cmd/database"
	models "chat/model"
	"log"
)

func main() {
	database.ConnectDB()

	if database.DB == nil {
		log.Fatal("❌ Database belum diinisialisasi")
	}

	err := database.DB.AutoMigrate(&models.User{}, &models.RoomChat{}, &models.RoomMember{}, &models.Chat{})
	if err != nil {
		log.Fatalf("❌ Gagal melakukan migrasi: %v", err)
	}

	log.Println("✅ Migrasi sukses!")
}
