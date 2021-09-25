package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/validators"
	"errors"
)

func CheckUserEmailExist(email string, filterIds []string) error {
	userModel := models.Users{}
	userList := userModel.UserInfosByEmailFilterIds(email, filterIds)
	if len(userList) != 0 {
		return errors.New(enums.CodeMessages(enums.UserEmailExist))
	}

	return nil
}

func UserCreate(userData *validators.UserRegister) error {
	userModel := &models.Users{
		Name:     userData.Name,
		Email:    userData.Email,
		Password: userData.Password,
	}

	addErr := userModel.UserAdd(userModel)

	return addErr
}
