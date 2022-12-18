package models

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"errors"
	"gorm.io/gorm"
	"strings"
)

type Services struct {
	ID       int64  `gorm:"column:id;primary_key"` //Service id
	ResID    string `gorm:"column:res_id"`         //ServiceResID
	Name     string `gorm:"column:name"`           //Service name
	Protocol int    `gorm:"column:protocol"`       //Protocol  1:HTTP  2:HTTPS  3:HTTP&HTTPS
	Enable   int    `gorm:"column:enable"`         //Service enable  1:on  2:off
	Release  int    `gorm:"column:release"`        //Service release status 1:unpublished  2:to be published  3:published
	ModelTime
}

// TableName sets the insert table name for this struct type
func (s *Services) TableName() string {
	return "oak_services"
}

var recursionTimesServices = 1

func (m *Services) ModelUniqueId() (string, error) {
	generateId, err := utils.IdGenerate(utils.IdTypeService)
	if err != nil {
		return "", err
	}

	result := packages.GetDb().Table(m.TableName()).Where("res_id = ?", generateId).Select("res_id").First(m)

	if result.RowsAffected == 0 {
		recursionTimesServices = 1
		return generateId, nil
	} else {
		if recursionTimesServices == utils.IdGenerateMaxTimes {
			recursionTimesServices = 1
			return "", errors.New(enums.CodeMessages(enums.IdConflict))
		}

		recursionTimesServices++
		id, err := m.ModelUniqueId()

		if err != nil {
			return "", err
		}

		return id, nil
	}
}

func (s *Services) ServiceInfoById(serviceId string) (Services, error) {

	var service Services
	err := packages.GetDb().Table(s.TableName()).Where("res_id = ?", serviceId).First(&service).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return Services{}, err
	}

	if err == gorm.ErrRecordNotFound {
		return Services{}, errors.New(enums.CodeMessages(enums.ServiceNull))
	}

	return service, nil
}

func (s *Services) ServiceAdd(serviceInfo *Services, serviceDomains []string) (string, error) {

	serviceId, err := s.ModelUniqueId()
	if err != nil {
		return serviceId, err
	}

	err = packages.GetDb().Transaction(func(tx *gorm.DB) error {
		serviceInfo.ResID = serviceId
		if serviceInfo.Name == "" {
			serviceInfo.Name = serviceId
		}

		err := tx.Create(&serviceInfo).Error

		if err != nil {
			packages.Log.Error("create services error")
			return err
		}

		serviceDomainsParam := []ServiceDomains{}
		for _, serviceDomain := range serviceDomains {
			domainId, err := (&ServiceDomains{}).ModelUniqueId()
			if err != nil {
				return err
			}

			serviceDomainsParam = append(serviceDomainsParam, ServiceDomains{
				Domain:       serviceDomain,
				ResID:        domainId,
				ServiceResID: serviceId,
			})
		}

		err = tx.Create(&serviceDomainsParam).Error

		if err != nil {
			packages.Log.Error("create services domain error")
			return err
		}

		return nil
	})
	return serviceId, err
}

func (s *Services) ServiceUpdate(serviceId string, serviceInfo *Services, serviceDomains []string) error {

	err := packages.GetDb().Transaction(func(tx *gorm.DB) error {
		err := tx.Table(s.TableName()).Where("res_id = ?", serviceId).Updates(serviceInfo).Error

		if err != nil {
			return nil
		}

		err = tx.Model(&ServiceDomains{}).Where("service_res_id = ?", serviceId).Delete(&ServiceDomains{}).Error

		if err != nil {
			return nil
		}

		serviceDomainsParam := []ServiceDomains{}
		for _, serviceDomain := range serviceDomains {
			domainId, err := (&ServiceDomains{}).ModelUniqueId()
			if err != nil {
				return err
			}

			serviceDomainsParam = append(serviceDomainsParam, ServiceDomains{
				Domain:       serviceDomain,
				ResID:        domainId,
				ServiceResID: serviceId,
			})
		}

		err = tx.Create(&serviceDomainsParam).Error

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Services) ServiceDelete(serviceId string) error {

	err := packages.GetDb().Transaction(func(tx *gorm.DB) error {

		err := tx.Table(s.TableName()).Where("res_id = ?", serviceId).Delete(Services{}).Error

		if err != nil {
			return err
		}

		err = tx.Model(&ServiceDomains{}).Where("service_res_id = ?", serviceId).Delete(ServiceDomains{}).Error

		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		packages.Log.Error("delete service error")
		return err
	}

	return nil
}

func (s *Services) ServiceList(serviceIds []string, param *validators.ServiceList) ([]Services, int, error) {

	tx := packages.GetDb().Table(s.TableName())

	if len(serviceIds) != 0 {
		tx.Where("res_id IN ?", serviceIds)
	}
	if param.Protocol != 0 {
		tx.Where("protocol = ?", param.Protocol)
	}
	if param.Enable != 0 {
		tx.Where("enable = ?", param.Enable)
	}
	if param.Release != 0 {
		tx.Where("release = ?", param.Release)
	}

	var total int
	err := ListCount(tx, &total)
	if err != nil {
		return []Services{}, 0, err
	}

	var list []Services
	tx = tx.Order("created_at desc")

	err = ListPaginate(tx, &list, &param.BaseListPage)

	if err != nil {
		return []Services{}, 0, err
	}
	return list, total, nil
}

func (s *Services) ServiceInfosLikeResIdName(resIdOrName string) ([]Services, error) {

	var service []Services
	resIdOrName = strings.TrimSpace(resIdOrName)
	if len(resIdOrName) == 0 {
		return []Services{}, nil
	}

	resIdOrName = "%" + resIdOrName + "%"
	err := packages.GetDb().Table(s.TableName()).
		Where("res_id LIKE ?", resIdOrName).
		Or("name LIKE ?", resIdOrName).
		Find(&service).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return service, nil
}

func (s *Services) ServiceUpdateColumns(serviceId string, params map[string]interface{}) error {

	err := packages.GetDb().Table(s.TableName()).Where("res_id = ?", serviceId).Updates(params).Error

	if err != nil {
		return err
	}

	return nil
}

func (s *Services) ServiceUpdateColumnsWithDB(tx *gorm.DB, serviceId string, params map[string]interface{}) error {

	err := tx.Table(s.TableName()).Where("res_id = ?", serviceId).Updates(params).Error

	if err != nil {
		return err
	}

	return nil
}
