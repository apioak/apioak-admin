package models

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"strings"
	"time"
)

type Certificates struct {
	ID          string    `gorm:"column:id;primary_key"` //Certificate id
	Sni         string    `gorm:"column:sni"`            //SNI
	Certificate string    `gorm:"column:certificate"`    //Certificate content
	PrivateKey  string    `gorm:"column:private_key"`    //Private key content
	IsEnable    int       `gorm:"column:is_enable"`      //Certificate enable  1:on  2:off
	ExpiredAt   time.Time `gorm:"column:expired_at"`     //Expiration time
	ModelTime
}

// TableName sets the insert table name for this struct type
func (c *Certificates) TableName() string {
	return "oak_certificates"
}

var certificatesId = ""

func (c *Certificates) CertificatesIdUnique(cIds map[string]string) (string, error) {
	if c.ID == "" {
		tmpID, err := utils.IdGenerate(utils.IdTypeCertificate)
		if err != nil {
			return "", err
		}
		c.ID = tmpID
	}

	result := packages.GetDb().
		Table(c.TableName()).
		Select("id").
		First(&c)

	mapId := cIds[c.ID]
	if (result.RowsAffected == 0) && (c.ID != mapId) {
		certificatesId = c.ID
		cIds[c.ID] = c.ID
		return certificatesId, nil
	} else {
		certId, certIdErr := utils.IdGenerate(utils.IdTypeCertificate)
		if certIdErr != nil {
			return "", certIdErr
		}
		c.ID = certId
		_, err := c.CertificatesIdUnique(cIds)
		if err != nil {
			return "", err
		}
	}

	return certificatesId, nil
}

func (c *Certificates) CertificatesAdd(certificatesData *Certificates) error {
	tpmIds := map[string]string{}
	certificatesId, certificatesIdUniqueErr := c.CertificatesIdUnique(tpmIds)
	if certificatesIdUniqueErr != nil {
		return certificatesIdUniqueErr
	}
	certificatesData.ID = certificatesId

	err := packages.GetDb().
		Table(c.TableName()).
		Create(certificatesData).Error

	return err
}

func (c *Certificates) CertificatesUpdate(id string, certificatesData *Certificates) error {
	updateError := packages.GetDb().
		Table(c.TableName()).
		Where("id = ?", id).
		Updates(certificatesData).Error

	return updateError
}

func (c *Certificates) CertificateInfoById(id string) Certificates {
	certificateInfo := Certificates{}
	packages.GetDb().
		Table(c.TableName()).
		Where("id = ?", id).
		First(&certificateInfo)

	return certificateInfo
}

func (c *Certificates) CertificateInfoBySni(sni string, filterId string) Certificates {
	certificateInfo := Certificates{}
	db := packages.GetDb().
		Table(c.TableName()).
		Where("sni = ?", sni)

	if len(filterId) != 0 {
		db = db.Where("id != ?", filterId)
	}

	db.First(&certificateInfo)

	return certificateInfo
}

func (c *Certificates) CertificateListPage(param *validators.CertificateList) (list []Certificates, total int, listError error) {
	tx := packages.GetDb().
		Table(c.TableName())

	if param.IsEnable != 0 {
		tx = tx.Where("is_enable = ?", param.IsEnable)
	}

	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		tx = tx.Where("sni LIKE ?", search)
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

func (c *Certificates) CertificateDelete(id string) error {
	deleteErr := packages.GetDb().
		Table(c.TableName()).
		Where("id = ?", id).
		Delete(c).Error

	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func (c *Certificates) CertificateSwitchEnable(id string, enable int) error {
	updateErr := packages.GetDb().
		Table(c.TableName()).
		Where("id = ?", id).
		Update("is_enable", enable).Error

	if updateErr != nil {
		return updateErr
	}

	return nil
}
