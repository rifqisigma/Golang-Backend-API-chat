package usecase

import (
	"chat/internal/repository"
	"chat/model"
	"chat/response"
	"chat/utils"
)

type RoomChatUseCase interface {
	GetGroupsByUserID(userID uint) ([]response.GetGroupByUserIdResponse, error)
	GetRoomChatByID(roomId uint) (*response.GetGroupByIdResponse, error)
	CreateRoom(userID uint, roomName string, roomDesc string) error
	DeleteRoom(roomID, adminID uint) error
	UpdateRoom(roomID, adminID uint, desc, name string) error

	//room member
	GetRoomMember(roomID uint) ([]response.RoomMemberResponse, error)
	Addmembers(roomID, adminID uint, targetIDS []uint) error
	DeleteMembersByAdmin(roomID, adminID uint, targetIDS []uint) error
	LeaveRoom(roomID, userID, targetID uint) error
}

type roomChatUseCase struct {
	roomChatRepo repository.RoomChatRepository
}

func NewRoomChatUseCase(roomChatRepo repository.RoomChatRepository) RoomChatUseCase {
	return &roomChatUseCase{roomChatRepo}
}

func (u *roomChatUseCase) GetRoomChatByID(roomId uint) (*response.GetGroupByIdResponse, error) {
	exist, err := u.roomChatRepo.IsRoomExist(roomId)
	if !exist {
		return nil, err
	}
	return u.roomChatRepo.GetRoomChatByID(roomId)
}

func (u *roomChatUseCase) GetGroupsByUserID(userID uint) ([]response.GetGroupByUserIdResponse, error) {
	return u.roomChatRepo.GetRoomChatsByUserID(userID)
}

func (u *roomChatUseCase) CreateRoom(userID uint, roomName string, roomDesc string) error {
	roomChat := &model.RoomChat{
		Name:      roomName,
		Desc:      roomDesc,
		CreatorID: userID,
	}

	return u.roomChatRepo.CreateRoomChat(roomChat, userID)
}

func (u *roomChatUseCase) DeleteRoom(roomID uint, adminID uint) error {
	if _, err := u.roomChatRepo.IsRoomExist(roomID); err != nil {
		return utils.ErrRoomNotFound
	}

	exist, _ := u.roomChatRepo.IsUserInRoom(roomID, adminID)
	if !exist {
		return utils.ErrUnauthorized
	}

	isAdmin, err := u.roomChatRepo.IsUserIsAdmin(roomID, adminID)
	if !isAdmin {
		return utils.ErrUnauthorized
	}
	if err != nil {
		return err
	}

	return u.roomChatRepo.DeleteRoom(roomID)
}
func (u *roomChatUseCase) UpdateRoom(roomID, adminID uint, desc, name string) error {
	existroom, err := u.roomChatRepo.IsRoomExist(roomID)
	if !existroom {
		return utils.ErrRoomNotFound
	}
	if err != nil {
		return err
	}

	exist, err := u.roomChatRepo.IsUserInRoom(roomID, adminID)
	if !exist {
		return utils.ErrUnauthorized
	}
	if err != nil {
		return err
	}

	isAdmin, err := u.roomChatRepo.IsUserIsAdmin(roomID, adminID)
	if !isAdmin {
		return utils.ErrUnauthorized
	}
	if err != nil {
		return err
	}

	return u.roomChatRepo.UpdateRoom(roomID, name, desc)
}

// room member
func (u *roomChatUseCase) GetRoomMember(roomID uint) ([]response.RoomMemberResponse, error) {
	if _, err := u.roomChatRepo.IsRoomExist(roomID); err != nil {
		return nil, utils.ErrRoomNotFound
	}

	return u.roomChatRepo.GetRoomMember(roomID)
}

func (u *roomChatUseCase) Addmembers(roomID, adminID uint, targetIDS []uint) error {
	roomexist, err := u.roomChatRepo.IsRoomExist(roomID)
	if !roomexist {
		return utils.ErrRoomNotFound
	}
	if err != nil {
		return err
	}

	exist, err := u.roomChatRepo.IsUserInRoom(roomID, adminID)
	if !exist {
		return utils.ErrUnauthorized
	}
	if err != nil {
		return err
	}

	isAdmin, err := u.roomChatRepo.IsUserIsAdmin(roomID, adminID)
	if !isAdmin {
		return utils.ErrUnauthorized
	}
	if err != nil {
		return err
	}

	if err := u.roomChatRepo.AddMembers(roomID, targetIDS, adminID); err != nil {
		return err
	}

	return nil
}

func (u *roomChatUseCase) DeleteMembersByAdmin(roomID, adminID uint, targetIDS []uint) error {
	existroom, err := u.roomChatRepo.IsRoomExist(roomID)
	if !existroom {
		return utils.ErrRoomNotFound
	}
	if err != nil {
		return err
	}

	exist, err := u.roomChatRepo.IsUserInRoom(roomID, adminID)
	if !exist {
		return utils.ErrUnauthorized
	}
	if err != nil {
		return err
	}

	isAdmin, err := u.roomChatRepo.IsUserIsAdmin(roomID, adminID)
	if !isAdmin {
		return utils.ErrUnauthorized
	}
	if err != nil {
		return err
	}

	return u.roomChatRepo.DeleteMembersByAdmin(roomID, targetIDS, adminID)

}

func (u *roomChatUseCase) LeaveRoom(roomID, userID, targetID uint) error {
	existroom, err := u.roomChatRepo.IsRoomExist(roomID)
	if !existroom {
		return utils.ErrRoomNotFound
	}
	if err != nil {
		return err
	}

	exist, err := u.roomChatRepo.IsUserInRoom(roomID, userID)
	if !exist {
		return utils.ErrUnauthorized
	}
	if err != nil {
		return err
	}

	return u.roomChatRepo.LeaveRoom(roomID, userID, targetID)
}
