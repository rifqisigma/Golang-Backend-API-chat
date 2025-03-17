package usecase

import (
	"chat/internal/repository"
	"chat/response"
	"chat/utils"
)

type ChatUsecase interface {
	CreateChat(message string, roomID, userID uint) (*response.ChatResponse, error)
	GetAllChatByRoomId(roomid uint) ([]response.ChatResponse, error)
	UpdateChat(roomID, userID uint, message string) error
	DeleteChat(roomID, userID uint) error
}

type chatUsecase struct {
	chatRepo     repository.ChatRepository
	roomChatRepo repository.RoomChatRepository
}

func NewChatUsecase(chatRepo repository.ChatRepository, roomChatRepo repository.RoomChatRepository) ChatUsecase {
	return &chatUsecase{
		chatRepo:     chatRepo,
		roomChatRepo: roomChatRepo,
	}
}

func (u *chatUsecase) CreateChat(message string, roomID, userID uint) (*response.ChatResponse, error) {
	exist, err := u.roomChatRepo.IsUserInRoom(roomID, userID)
	if !exist {
		return nil, utils.ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}

	existroom, err := u.roomChatRepo.IsRoomExist(roomID)
	if !existroom {
		return nil, utils.ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}

	response, err := u.chatRepo.CreateChat(message, roomID, userID)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (u *chatUsecase) GetAllChatByRoomId(roomid uint) ([]response.ChatResponse, error) {
	exist, err := u.roomChatRepo.IsRoomExist(roomid)
	if !exist {
		return nil, utils.ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}

	return u.chatRepo.GetAllChatByRoomID(roomid)
}

func (u *chatUsecase) UpdateChat(roomID, userID uint, message string) error {
	exist, err := u.roomChatRepo.IsUserInRoom(roomID, userID)
	if !exist {
		return utils.ErrUnauthorized
	}
	if err != nil {
		return err
	}

	existroom, err := u.roomChatRepo.IsRoomExist(roomID)
	if !existroom {
		return utils.ErrRoomNotFound
	}
	if err != nil {
		return err
	}

	return u.chatRepo.UpdateChat(roomID, userID, message)
}

func (u *chatUsecase) DeleteChat(roomID, userID uint) error {
	exist, err := u.roomChatRepo.IsUserInRoom(roomID, userID)
	if !exist {
		return utils.ErrUnauthorized
	}
	if err != nil {
		return err
	}

	existroom, err := u.roomChatRepo.IsRoomExist(roomID)
	if !existroom {
		return utils.ErrRoomNotFound
	}
	if err != nil {
		return err
	}
	return u.chatRepo.DeleteChat(roomID, userID)

}
