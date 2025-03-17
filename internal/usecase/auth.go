package usecase

import (
	"chat/internal/repository"
	"chat/model"
	"fmt"

	"chat/utils"
	"errors"
)

type UserUsecase interface {
	Register(username, email, password string) error
	Login(email, password string) (*model.User, error)
	UpdateUser(firstEmail, username, email string) error
	DeleteUser(email string) error
	ValidateUser(email string) error
	ResendLinkVerif(email string) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo}
}

func (u *userUsecase) Register(username, email, password string) error {

	if !utils.IsValidEmail(email) {
		return errors.New("invalid email format")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	err = u.userRepo.CreateUser(username, email, hashedPassword)
	if err != nil {
		return utils.ErrInternal
	}

	tokenJWT, err := utils.GenerateJWTVerification(email)
	if err != nil {
		return err
	}

	utils.SendEmail(email, tokenJWT)

	return nil
}

func (u *userUsecase) Login(email, password string) (*model.User, error) {
	user, err := u.userRepo.Login(email)
	if err != nil {
		return nil, errors.New("failed to login")
	}

	if !utils.CheckPassword(user.Password, password) {
		return nil, errors.New("invalid password")
	}
	return user, nil
}

func (u *userUsecase) UpdateUser(firstEmail, username, email string) error {
	exist, err := u.userRepo.IsUserExist(email)
	if !exist {
		return utils.ErrUserNotFound
	}
	if err != nil {
		return fmt.Errorf("err")
	}
	return u.userRepo.UpdateUser(firstEmail, username, email)
}

func (u *userUsecase) DeleteUser(email string) error {
	exist, err := u.userRepo.IsUserExist(email)
	if !exist {
		return utils.ErrUserNotFound
	}
	if err != nil {
		return fmt.Errorf("err")
	}

	return u.userRepo.DeleteUser(email)
}

func (u *userUsecase) ValidateUser(email string) error {
	return u.userRepo.ValidateUser(email)
}

func (u *userUsecase) ResendLinkVerif(email string) error {
	if !utils.IsValidEmail(email) {
		return utils.ErrBadRequest
	}

	exist, err := u.userRepo.IsUserExist(email)
	if !exist {
		return utils.ErrUserNotFound
	}
	if err != nil {
		return err
	}

	tokenJWT, err := utils.GenerateJWTVerification(email)
	if err != nil {
		return err
	}

	utils.SendEmail(email, tokenJWT)

	return nil
}
