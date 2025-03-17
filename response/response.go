package response

import "time"

type GetGroupByUserIdResponse struct {
	ID   uint   `json:"room_id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type GetGroupByIdResponse struct {
	ID         uint                 `json:"id"`
	Name       string               `json:"name"`
	Desc       string               `json:"desc"`
	CreatorID  uint                 `json:"creator_id"`
	RoomMember []RoomMemberResponse `json:"room_members"`
}

type RoomMemberResponse struct {
	ID     uint   `json:"room_member_id"`
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}

type ChatResponse struct {
	RoomID   uint      `json:"room_id"`
	SenderID uint      `json:"user_id"`
	Message  string    `json:"message"`
	Time     time.Time `json:"time"`
}
