package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"errors"
	"fmt"
	"time"
)

func CheckUserEmailExist(email string, filterIds []string) error {
	userModel := models.Users{}
	userList := userModel.UserInfosByEmailFilterIds(email, filterIds)
	if len(userList) != 0 {
		return errors.New(enums.CodeMessages(enums.UserEmailExist))
	}

	return nil
}

func CheckUserAndPassword(email string, password string) error {
	userModel := models.Users{}
	userInfo := userModel.UserInfoByEmail(email)
	if userInfo.Email != email {
		return errors.New(enums.CodeMessages(enums.UserNull))
	}

	if utils.Md5(utils.Md5(password)) != userInfo.Password {
		return errors.New(enums.CodeMessages(enums.UserPasswordError))
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

func UserLogin(email string) (string, error) {
	token, tokenErr := utils.GenToken(email)
	if tokenErr != nil {
		return "", errors.New(enums.CodeMessages(enums.UserLoggingInError))
	}

	emailExpires, _ := time.ParseDuration(fmt.Sprintf("+%dm", packages.Token.TokenExpire))
	utils.SetTokenExpire(token, time.Now().Add(emailExpires).Unix())

	return token, nil
}

func UserLogout(token string) (bool, error) {
	email, err := utils.ParseToken(token)
	if err != nil {
		return false, errors.New(enums.CodeMessages(enums.UserTokenError))
	}

	allEmailExpireTimes := *utils.GetTokenExpire()
	delete(allEmailExpireTimes, email)

	return true, nil
}

func UserLoginRefresh(token string) (bool, error) {
	_, err := utils.ParseToken(token)
	if err != nil {
		return false, errors.New(enums.CodeMessages(enums.UserTokenError))
	}

	emailExpires, _ := time.ParseDuration(fmt.Sprintf("+%dm", packages.Token.TokenExpire))
	utils.SetTokenExpire(token, time.Now().Add(emailExpires).Unix())

	return true, nil
}

func CheckUserLoginStatus(token string) (bool, error) {
	email, err := utils.ParseToken(token)
	if err != nil {
		return false, errors.New(enums.CodeMessages(enums.UserTokenError))
	}

	allEmailExpireTimes := *utils.GetTokenExpire()
	emailExpireTime, emailExpireTimeExist := allEmailExpireTimes[email]

	if emailExpireTimeExist == false {

		return false, errors.New(enums.CodeMessages(enums.UserNoLoggingIn))

	} else {
		if emailExpireTime.Token != token {
			return false, errors.New(enums.CodeMessages(enums.UserTokenError))
		}

		if emailExpireTime.Expire < time.Now().Unix() {
			return false, errors.New(enums.CodeMessages(enums.UserLoggingInExpire))
		}
	}

	return true, nil
}
