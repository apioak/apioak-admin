package models

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"errors"
	"gorm.io/gorm"
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

func (m *UserTokens) ModelUniqueId() (generateId string, err error) {
	generateId, err = utils.IdGenerate(utils.IdTypeUserToken)
	if err != nil {
		return
	}

	err = packages.GetDb().
		Table(m.TableName()).
		Where("res_id = ?", generateId).
		Select("res_id").
		First(m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	if err == nil {
		recursionTimesServices = 1
		return
	} else {
		if recursionTimesServices == utils.IdGenerateMaxTimes {
			recursionTimesServices = 1
			err = errors.New(enums.CodeMessages(enums.IdConflict))
			return
		}

		recursionTimesServices++
		generateId, err = m.ModelUniqueId()
		if err != nil {
			return
		}

		return
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
