package repository

import (
	"chat/model"
	"chat/response"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type RoomChatRepository interface {
	GetRoomChatsByUserID(userID uint) ([]response.GetGroupByUserIdResponse, error)
	GetRoomChatByID(roomID uint) (*response.GetGroupByIdResponse, error)
	CreateRoomChat(roomChat *model.RoomChat, userID uint) error
	UpdateRoom(roomID uint, name, desc string) error
	DeleteRoom(roomID uint) error
	IsRoomExist(roomID uint) (bool, error)

	//for room member
	GetRoomMember(roomID uint) ([]response.RoomMemberResponse, error)
	AddMembers(roomID uint, targetIDS []uint, adminID uint) error
	IsUserInRoom(roomID uint, userID uint) (bool, error)
	IsUserIsAdmin(roomID uint, userID uint) (bool, error)
	DeleteMembersByAdmin(roomID uint, targetIDS []uint, userID uint) error
	LeaveRoom(roomID uint, userID uint, targetID uint) error
}

type roomChatRepository struct {
	db *gorm.DB
}

func NewRoomChatRepository(db *gorm.DB) RoomChatRepository {
	return &roomChatRepository{db}
}

func (r *roomChatRepository) GetRoomChatsByUserID(userID uint) ([]response.GetGroupByUserIdResponse, error) {

	var rooms []response.GetGroupByUserIdResponse
	err := r.db.
		Model(&model.RoomChat{}).
		Select("room_chats.id, room_chats.name, room_chats.desc").
		Joins("JOIN room_members rm ON rm.room_id = room_chats.id").
		Where("rm.user_id = ?", userID).
		Find(&rooms).Error

	return rooms, err
}

func (r *roomChatRepository) GetRoomChatByID(roomID uint) (*response.GetGroupByIdResponse, error) {
	var room model.RoomChat
	if err := r.db.Model(&room).
		Select("room_chats.id, room_chats.name, room_chats.desc, room_chats.creator_id").
		Where("id = ?", roomID).First(&room).Error; err != nil {
		return nil, err
	}
	responses := response.GetGroupByIdResponse{
		ID:        room.ID,
		Name:      room.Name,
		Desc:      room.Desc,
		CreatorID: room.CreatorID,
	}

	var roomMembers []response.RoomMemberResponse
	if err := r.db.Model(&model.RoomMember{}).
		Select("room_members.id, room_members.user_id, room_members.role").
		Where("room_id = ?", roomID).Find(&roomMembers).Error; err != nil {
		return nil, err
	}

	responses.RoomMember = roomMembers

	return &responses, nil
}

func (r *roomChatRepository) CreateRoomChat(roomChat *model.RoomChat, userID uint) error {
	tx := r.db.Begin()

	if err := tx.Create(roomChat).Error; err != nil {
		tx.Rollback()
		return err
	}

	roomMember := &model.RoomMember{
		RoomID: roomChat.ID,
		Role:   "admin",
		UserID: &userID,
	}

	if err := tx.Create(&roomMember).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *roomChatRepository) UpdateRoom(roomID uint, name, desc string) error {
	result := r.db.Model(&model.RoomChat{}).Where("id= ? ", roomID).Updates(map[string]interface{}{
		"name": name,
		"desc": desc,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no chat found to update")
	}
	return nil
}

func (r *roomChatRepository) DeleteRoom(roomID uint) error {
	if err := r.db.Where("id = ?", roomID).Delete(&model.RoomChat{}).Error; err != nil {
		return fmt.Errorf("failed to delete room members: %v", err)
	}
	return nil
}

func (r *roomChatRepository) IsUserInRoom(roomID uint, userID uint) (bool, error) {
	var room model.RoomMember
	if err := r.db.Where("room_id = ? AND user_id = ?", roomID, userID).First(&room).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, err
		}
		return false, err
	}
	return true, nil
}

func (r *roomChatRepository) IsUserIsAdmin(roomID uint, userID uint) (bool, error) {
	var count int64
	if err := r.db.Model(&model.RoomMember{}).Where("room_id = ? AND user_id = ? AND role = ?", roomID, userID, "admin").Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *roomChatRepository) IsRoomExist(roomID uint) (bool, error) {

	if err := r.db.First(&model.RoomChat{}, roomID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
		return false, err
	}
	return true, nil
}

// room member
func (r *roomChatRepository) GetRoomMember(roomID uint) ([]response.RoomMemberResponse, error) {
	var response []response.RoomMemberResponse
	if err := r.db.Model(&model.RoomMember{}).Where("room_id = ?", roomID).Find(&response).Error; err != nil {
		return nil, err
	}
	return response, nil
}

func (r *roomChatRepository) AddMembers(roomID uint, targetIDs []uint, adminID uint) error {
	var isExist int64
	r.db.Model(&model.RoomMember{}).
		Where("room_id = ? AND user_id IN ?", roomID, targetIDs).
		Count(&isExist)

	if isExist > 0 {
		return fmt.Errorf("some users are already in the group, operation aborted")
	}

	var newMembers []model.RoomMember
	for _, userID := range targetIDs {
		newMembers = append(newMembers, model.RoomMember{
			RoomID: roomID,
			UserID: &userID,
			Role:   "member",
		})
	}

	if err := r.db.Create(&newMembers).Error; err != nil {
		return err
	}

	return nil
}

func (r *roomChatRepository) DeleteMembersByAdmin(roomID uint, targetIDS []uint, userID uint) error {
	var count int64
	r.db.Model(&model.RoomMember{}).
		Where("room_id = ? AND user_id IN ?", roomID, targetIDS).
		Count(&count)

	if count != int64(len(targetIDS)) {
		return fmt.Errorf("some users are not in the room, operation aborted")
	}

	if err := r.db.Where("room_id = ? AND user_id IN ?", roomID, targetIDS).Delete(&model.RoomMember{}).Error; err != nil {
		return err
	}
	return nil

}

func (r *roomChatRepository) LeaveRoom(roomID uint, userID uint, targetID uint) error {
	return r.db.Where("room_id = ? AND user_id = ?", roomID, targetID).Delete(&model.RoomMember{}).Error
}
