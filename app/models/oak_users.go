package models

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
)

type Users struct {
	ID       string `gorm:"column:id;primary_key"` //User iD
	Name     string `gorm:"column:name"`           //User name
	Password string `gorm:"column:password"`       //Password
	Email    string `gorm:"column:email"`          //Email
	ModelTime
}

// TableName sets the insert table name for this struct type
func (u *Users) TableName() string {
	return "oak_users"
}

var userId = ""

func (u *Users) UserIdUnique(uIds map[string]string) (string, error) {
	if u.ID == "" {
		tmpID, err := utils.IdGenerate(utils.IdTypeUser)
		if err != nil {
			return "", err
		}
		u.ID = tmpID
	}

	result := packages.GetDb().
		Table(u.TableName()).
		Select("id").
		First(&u)

	mapId := uIds[u.ID]
	if (result.RowsAffected == 0) && (u.ID != mapId) {
		userId = u.ID
		uIds[u.ID] = u.ID
		return userId, nil
	} else {
		usId, certIdErr := utils.IdGenerate(utils.IdTypeUser)
		if certIdErr != nil {
			return "", certIdErr
		}
		u.ID = usId
		_, err := u.UserIdUnique(uIds)
		if err != nil {
			return "", err
		}
	}

	return userId, nil
}

func (u *Users) UserInfosByEmailFilterIds(email string, filterIds []string) []Users {
	userInfos := make([]Users, 0)
	db := packages.GetDb().
		Table(u.TableName()).
		Where("email = ?", email)

	if len(filterIds) != 0 {
		db = db.Where("id NOT IN ?", filterIds)
	}

	db.Find(&userInfos)

	return userInfos
}

func (u *Users) UserAdd(userData *Users) error {
	tpmIds := map[string]string{}
	userId, userIdUniqueErr := u.UserIdUnique(tpmIds)
	if userIdUniqueErr != nil {
		return userIdUniqueErr
	}
	userData.ID = userId
	userData.Password = utils.Md5(utils.Md5(userData.Password))

	err := packages.GetDb().
		Table(u.TableName()).
		Create(userData).Error

	return err
}
