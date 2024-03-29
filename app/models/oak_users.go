package models

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"errors"
	"gorm.io/gorm"
)

type Users struct {
	ID       int    `gorm:"column:id;primary_key"` // primary key
	ResID    string `gorm:"column:res_id"`         // User iD
	Name     string `gorm:"column:name"`           // User name
	Password string `gorm:"column:password"`       // Password
	Email    string `gorm:"column:email"`          // Email
	ModelTime
}

// TableName sets the insert table name for this struct type
func (u *Users) TableName() string {
	return "oak_users"
}

var recursionTimesUsers = 1

func (m *Users) ModelUniqueId() (generateId string, err error) {
	generateId, err = utils.IdGenerate(utils.IdTypeUser)
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
	userId, userIdUniqueErr := u.ModelUniqueId()
	if userIdUniqueErr != nil {
		return userIdUniqueErr
	}
	userData.ResID = userId
	userData.Password = utils.Md5(utils.Md5(userData.Password))

	err := packages.GetDb().
		Table(u.TableName()).
		Create(userData).Error

	return err
}

func (u *Users) UserInfoByEmail(email string) Users {
	userInfo := Users{}
	packages.GetDb().
		Table(u.TableName()).
		Where("email = ?", email).
		First(&userInfo)

	return userInfo
}
