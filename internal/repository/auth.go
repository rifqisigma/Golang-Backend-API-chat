package repository

import (
	"chat/model"
	"errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(username, email, hashedpassword string) error
	Login(email string) (*model.User, error)
	UpdateUser(firstEmail, username, email string) error
	DeleteUser(email string) error
	ValidateUser(email string) error
	IsValidate(email string) (bool, error)
	IsUserExist(email string) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateUser(username, email, hashedpassword string) error {
	newUser := model.User{
		Username: username,
		Email:    email,
		Password: hashedpassword,
	}

	if err := r.db.Create(&newUser).Error; err != nil {
		return err
	}

	return nil

}

func (r *userRepository) Login(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Limit(1).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(firstEmail, username, email string) error {
	user := model.User{
		Email:    email,
		Username: username,
	}
	if err := r.db.Model(&model.User{}).Where("email = ?", firstEmail).Updates(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) DeleteUser(email string) error {
	tx := r.db.Begin()

	var user model.User
	if err := tx.Where("email = ?", email).First(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("sender_id", user.ID).Delete(&model.Chat{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", user.ID).Delete(&model.RoomMember{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("creator_id = ?", user.ID).Delete(&model.RoomChat{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", user.ID).Delete(&model.User{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error

}

func (r *userRepository) ValidateUser(email string) error {
	result := r.db.Model(&model.User{}).
		Where("email = ?", email).
		Update("is_verified", true) // Bisa pakai Update tunggal

	if result.Error != nil {
		return result.Error // Kembalikan error jika query gagal
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found or already verified") // Cek apakah ada data yang berubah
	}

	return nil
}

func (r *userRepository) IsValidate(email string) (bool, error) {
	err := r.db.Model(&model.User{}).Where("email = ? AND is_verified", email, true)
	if err != nil {
		return false, err.Error
	}
	return true, nil
}

func (r *userRepository) IsUserExist(email string) (bool, error) {
	var user model.User

	// Cek apakah ada user dengan email tersebut
	err := r.db.Model(&model.User{}).Where("email = ?", email).First(&user).Error

	if err != nil {
		// Jika error karena tidak ditemukan, berarti user tidak ada
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		// Jika error lain, kembalikan error
		return false, err
	}

	// Jika tidak error, berarti user ditemukan
	return true, nil
}
