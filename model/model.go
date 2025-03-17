package model

import (
	"time"
)

// User menyimpan data pengguna
type User struct {
	ID         uint      `gorm:"primaryKey"`
	Username   string    `gorm:"unique;not null"`
	Email      string    `gorm:"unique;not null"`
	Password   string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	IsVerified bool      `gorm:"default:false"`
}

type RoomChat struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Desc      string
	CreatorID uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Relasi
	Creator User `gorm:"foreignKey:CreatorID"`

	// Relasi for delete cascade
	RoomMembers []RoomMember `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE"`
	Chats       []Chat       `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE"`
}

type RoomMember struct {
	ID     uint   `gorm:"primaryKey"`
	RoomID uint   `gorm:"not null;uniqueIndex:idx_room_user"`
	UserID *uint  `gorm:"uniqueIndex:idx_room_user;constraint:OnDelete:SET NULL"`
	Role   string `gorm:"not null"`

	// Relasi
	Room RoomChat `gorm:"foreignKey:RoomID"`
	User User     `gorm:"foreignKey:UserID"`
}

type Chat struct {
	ID        uint      `gorm:"primaryKey"`
	RoomID    uint      `gorm:"not null;index"`
	SenderID  *uint     `gorm:"constraint:OnDelete:SET NULL"`
	Message   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Relasi
	Room   RoomChat `gorm:"foreignKey:RoomID"`
	Sender User     `gorm:"foreignKey:SenderID"`
}
