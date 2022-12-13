package models

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"errors"
	"strings"
	"time"
)

type Certificates struct {
	ID          string    `gorm:"column:id;primary_key"` //Certificate id
	Sni         string    `gorm:"column:sni"`            //SNI
	Certificate string    `gorm:"column:certificate"`    //Certificate content
	PrivateKey  string    `gorm:"column:private_key"`    //Private key content
	Enable      int       `gorm:"column:enable"`         //Certificate enable  1:on  2:off
	Release     int       `gorm:"column:release"`        //Certificates release status 1:unpublished  2:to be published  3:published
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
		Where("id = ?", generateId).
		Select("id").
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
		id, err := m.ModelUniqueId()

		if err != nil {
			return "", err
		}

		return id, nil
	}
}

func (c *Certificates) CertificatesAdd(certificatesData *Certificates) (string, error) {
	certificatesId, certificatesIdUniqueErr := c.ModelUniqueId()
	if certificatesIdUniqueErr != nil {
		return certificatesId, certificatesIdUniqueErr
	}
	certificatesData.ID = certificatesId

	err := packages.GetDb().
		Table(c.TableName()).
		Create(certificatesData).Error

	return certificatesId, err
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
	if param.ReleaseStatus != 0 {
		tx = tx.Where("release_status = ?", param.ReleaseStatus)
	}

	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		tx = tx.Where(
			packages.GetDb().Table(c.TableName()).
				Where("sni LIKE ?", search).
				Or("id LIKE ?", search))
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
	certificateInfo := c.CertificateInfoById(id)
	releaseStatus := certificateInfo.Release
	if certificateInfo.Release == utils.ReleaseStatusY {
		releaseStatus = utils.ReleaseStatusT
	}

	updateErr := packages.GetDb().
		Table(c.TableName()).
		Where("id = ?", id).
		Updates(Certificates{
			Enable:  enable,
			Release: releaseStatus}).Error

	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (c *Certificates) CertificateSwitchRelease(id string, release int) error {
	updateErr := packages.GetDb().
		Table(c.TableName()).
		Where("id = ?", id).
		Update("release_status", release).Error

	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (c *Certificates) CertificateInfoByDomainSniInfos(domains []string) []Certificates {
	certificateInfos := make([]Certificates, 0)
	if len(domains) == 0 {
		return certificateInfos
	}

	packages.GetDb().
		Table(c.TableName()).
		Where("sni In ?", domains).
		Find(&certificateInfos)

	return certificateInfos
}
