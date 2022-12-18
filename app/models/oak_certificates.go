package models

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"errors"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Certificates struct {
	ID          int       `gorm:"column:id;primary_key"` //Certificate id
	ResID       string    `gorm:"column:res_id"`         //SNI
	Sni         string    `gorm:"column:sni"`            //SNI
	Certificate string    `gorm:"column:certificate"`    //Certificate content
	PrivateKey  string    `gorm:"column:private_key"`    //Private key content
	Enable      int       `gorm:"column:enable"`         //Certificate enable  1:on  2:off
	ExpiredAt   time.Time `gorm:"column:expired_at"`     //Expiration time
	ModelTime
}

// TableName sets the insert table name for this struct type
func (c *Certificates) TableName() string {
	return "oak_certificates"
}

var recursionTimesCertificates = 1

func (m *Certificates) ModelUniqueId() (string, error) {
	generateId, generateIdErr := utils.IdGenerate(utils.IdTypeCertificate)
	if generateIdErr != nil {
		return "", generateIdErr
	}

	result := packages.GetDb().
		Table(m.TableName()).
		Where("res_id = ?", generateId).
		Select("res_id").
		First(m)

	if result.RowsAffected == 0 {
		recursionTimesCertificates = 1
		return generateId, nil
	} else {
		if recursionTimesCertificates == utils.IdGenerateMaxTimes {
			recursionTimesCertificates = 1
			return "", errors.New(enums.CodeMessages(enums.IdConflict))
		}

		recursionTimesCertificates++
		resID, err := m.ModelUniqueId()

		if err != nil {
			return "", err
		}

		return resID, nil
	}
}

func (c *Certificates) CertificatesAdd(tx *gorm.DB, certificatesData *Certificates) (string, error) {
	certificatesId, err := c.ModelUniqueId()
	if err != nil {
		return certificatesId, err
	}

	certificatesData.ResID = certificatesId

	err = tx.Table(c.TableName()).Create(certificatesData).Error

	return certificatesId, err
}

func (c *Certificates) CertificatesUpdate(tx *gorm.DB, resID string, certificatesData *Certificates) error {

	return tx.Table(c.TableName()).Where("res_id = ?", resID).Updates(certificatesData).Error

}

func (c *Certificates) CertificateInfoById(resID string) (Certificates, error) {
	certificateInfo := Certificates{}
	err := packages.GetDb().Table(c.TableName()).Where("res_id = ?", resID).First(&certificateInfo).Error

	if err != nil {
		return certificateInfo, err
	}

	return certificateInfo, nil
}

func (c *Certificates) EnableCertificateInfoBySni(sni string, filterId string) (Certificates, error) {
	certificateInfo := Certificates{}
	db := packages.GetDb().Table(c.TableName()).Where("sni = ?", sni)

	if filterId != "" {
		db = db.Where("res_id != ?", filterId)
	}

	err := db.Where("enable = ?", utils.EnableOn).First(&certificateInfo).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return certificateInfo, err
	}
	return certificateInfo, nil
}

func (c *Certificates) CertificateListPage(param *validators.CertificateList) (list []Certificates, total int, listError error) {
	tx := packages.GetDb().
		Table(c.TableName())

	if param.Enable != 0 {
		tx = tx.Where("enable = ?", param.Enable)
	}

	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		tx = tx.Where(
			packages.GetDb().Table(c.TableName()).
				Where("sni LIKE ?", search).
				Or("res_id LIKE ?", search))
	}

	countError := ListCount(tx, &total)
	if countError != nil {
		listError = countError
		return
	}

	tx = tx.Order("expired_at ASC")
	listError = ListPaginate(tx, &list, &param.BaseListPage)
	return
}

func (c *Certificates) CertificateDelete(tx *gorm.DB, resID string) error {
	err := tx.Model(&Certificates{}).
		Where("res_id = ?", resID).
		Delete(&Certificates{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (c *Certificates) CertificateSwitchEnable(tx *gorm.DB, resID string, enable int) error {

	updateParam := Certificates{
		Enable: enable,
	}
	err := tx.Table(c.TableName()).Where("res_id = ?", resID).Updates(updateParam).Error

	if err != nil {
		return err
	}

	return nil
}

func (c *Certificates) CertificateInfoByDomainSniInfos(domains []string) []Certificates {
	certificateInfos := []Certificates{}
	if len(domains) == 0 {
		return certificateInfos
	}

	packages.GetDb().Table(c.TableName()).Where("sni In ?", domains).Find(&certificateInfos)

	return certificateInfos
}
