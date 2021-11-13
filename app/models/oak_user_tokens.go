package models

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"time"
)

type UserTokens struct {
	ID        string    `gorm:"column:id;primary_key"` //User tokenID
	Token     string    `gorm:"column:token"`          //Token
	UserEmail string    `gorm:"column:user_email"`     //Email
	CreatedAt time.Time `gorm:"column:created_at"`     //Creation time
	UpdatedAt time.Time `gorm:"column:updated_at"`     //Update time
	ExpiredAt time.Time `gorm:"column:expired_at"`     //Expired time
}

// TableName sets the insert table name for this struct type
func (u *UserTokens) TableName() string {
	return "oak_user_tokens"
}

var userTokenId = ""

func (u *UserTokens) UserTokenIdUnique(utIds map[string]string) (string, error) {
	if u.ID == "" {
		tmpID, err := utils.IdGenerate(utils.IdTypeUserToken)
		if err != nil {
			return "", err
		}
		u.ID = tmpID
	}

	result := packages.GetDb().
		Table(u.TableName()).
		Select("id").
		First(&u)

	mapId := utIds[u.ID]
	if (result.RowsAffected == 0) && (u.ID != mapId) {
		userTokenId = u.ID
		utIds[u.ID] = u.ID
		return userTokenId, nil
	} else {
		ustId, certIdErr := utils.IdGenerate(utils.IdTypeUserToken)
		if certIdErr != nil {
			return "", certIdErr
		}
		u.ID = ustId
		_, err := u.UserTokenIdUnique(utIds)
		if err != nil {
			return "", err
		}
	}

	return userTokenId, nil
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
		tpmIds := map[string]string{}
		userTokenId, userIdUniqueErr := u.UserTokenIdUnique(tpmIds)
		if userIdUniqueErr != nil {
			return userIdUniqueErr
		}
		updateData := &UserTokens{}
		updateData.ID = userTokenId
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
