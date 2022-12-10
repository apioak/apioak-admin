package models

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"errors"
	"time"
)

type UserTokens struct {
	ID        int       `gorm:"column:id;primary_key"` //primary key
	ResID     string    `gorm:"column:res_id"`         //User tokenID
	Token     string    `gorm:"column:token"`          //Token
	UserEmail string    `gorm:"column:user_email"`     //Email
	ExpiredAt time.Time `gorm:"column:expired_at"`     //Expired time
	ModelTime
}

// TableName sets the insert table name for this struct type
func (u *UserTokens) TableName() string {
	return "oak_user_tokens"
}

var recursionTimesUserTokens = 1

func (m *UserTokens) ModelUniqueId() (string, error) {
	generateId, generateIdErr := utils.IdGenerate(utils.IdTypeUserToken)
	if generateIdErr != nil {
		return "", generateIdErr
	}

	result := packages.GetDb().
		Table(m.TableName()).
		Where("res_id = ?", generateId).
		Select("res_id").
		First(m)

	if result.RowsAffected == 0 {
		recursionTimesUserTokens = 1
		return generateId, nil
	} else {
		if recursionTimesUserTokens == utils.IdGenerateMaxTimes {
			recursionTimesUserTokens = 1
			return "", errors.New(enums.CodeMessages(enums.IdConflict))
		}

		recursionTimesUserTokens++
		id, err := m.ModelUniqueId()

		if err != nil {
			return "", err
		}

		return id, nil
	}
}

func (u *UserTokens) SetTokenExpire(email string, token string, expiredTime time.Time) error {
	packages.GetDb().
		Table(u.TableName()).
		Where("user_email = ?", email).
		First(u)

	if len(u.UserEmail) != 0 && u.UserEmail == email {
		updateData := UserTokens{}
		updateData.Token = token
		updateData.ExpiredAt = expiredTime

		packages.GetDb().
			Table(u.TableName()).
			Where("user_email = ?", email).
			Updates(updateData)

	} else {
		userTokenId, userIdUniqueErr := u.ModelUniqueId()
		if userIdUniqueErr != nil {
			return userIdUniqueErr
		}
		updateData := &UserTokens{}
		updateData.ResID = userTokenId
		updateData.UserEmail = email
		updateData.Token = token
		updateData.ExpiredAt = expiredTime

		packages.GetDb().
			Table(u.TableName()).
			Create(updateData)
	}

	return nil
}

func (u *UserTokens) GetTokenExpireByEmail(email string) *UserTokens {
	userTokenExpireInfo := UserTokens{}

	packages.GetDb().
		Table(u.TableName()).
		Where("user_email = ?", email).
		First(&userTokenExpireInfo)

	return &userTokenExpireInfo
}

func (u *UserTokens) DelTokenExpireByEmail(email string) error {
	delErr := packages.GetDb().
		Table(u.TableName()).
		Where("user_email = ?", email).
		Delete(u)

	return delErr.Error
}
