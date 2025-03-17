package repository

import (
	"chat/model"
	"chat/response"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ChatRepository interface {
	GetAllChatByRoomID(roomID uint) ([]response.ChatResponse, error)
	CreateChat(message string, roomID, userID uint) (*response.ChatResponse, error)
	DeleteChat(roomID, userID uint) error
	UpdateChat(roomID, userID uint, message string) error
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db}
}

func (r *chatRepository) GetAllChatByRoomID(roomID uint) ([]response.ChatResponse, error) {
	var response []response.ChatResponse
	err := r.db.Table("chats").
		Select("room_id,sender_id,message,created_at").
		Where("room_id = ?", roomID).
		Order("created_at ASC").
		Find(&response).Error

	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *chatRepository) CreateChat(message string, roomID, userID uint) (*response.ChatResponse, error) {
	chat := &model.Chat{
		RoomID:    roomID,
		SenderID:  &userID,
		Message:   message,
		CreatedAt: time.Now(),
	}

	if err := r.db.Create(chat).Error; err != nil {
		return nil, err
	}
	response := &response.ChatResponse{
		RoomID:   roomID,
		SenderID: userID,
		Message:  message,
		Time:     time.Now(),
	}
	return response, nil
}

func (r *chatRepository) DeleteChat(roomID, userID uint) error {
	result := r.db.Model(&model.Chat{}).Where("room_id = ? AND sender_id = ?", roomID, userID).Delete(&model.Chat{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no chat found to delete")
	}
	return nil
}

func (r *chatRepository) UpdateChat(roomID, userID uint, message string) error {
	result := r.db.Model(&model.Chat{}).Where("room_id = ? AND sender_id = ?", roomID, userID).Updates(map[string]interface{}{
		"message": message,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no chat found to update")
	}
	return nil
}
