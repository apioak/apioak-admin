package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/packages"
	"apioak-admin/app/rpc"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"errors"
	"gorm.io/gorm"
	"sync"
)

type CertificateService struct {
}

var (
	certificateService *CertificateService
	certificateOnce    sync.Once
)

func NewCertificateService() *CertificateService {

	certificateOnce.Do(func() {
		certificateService = &CertificateService{}
	})

	return certificateService
}

func syncDataSideCertificate(tx *gorm.DB, new *models.Certificates, filterID string) error {

	existSniCertificate, err := (&models.Certificates{}).EnableCertificateInfoBySni(new.Sni, filterID)
	if err != nil {
		return err
	}

	tmpDeleteRes := false
	// 相同域名的证书已启用时需要将已启用的证书关闭，并同步至数据面
	if existSniCertificate.Enable == utils.EnableOn {
		// 修改控制面旧证书启用状态
		err := (&models.Certificates{}).CertificateSwitchEnable(tx, existSniCertificate.ResID, utils.EnableOff)

		if err != nil {
			return nil
		}
		// 同步数据面旧证书启用状态
		err = rpc.NewApiOak().CertificateDelete(existSniCertificate.ResID)
		if err != nil {
			return nil
		}
		tmpDeleteRes = true
	}

	// 新增数据面证书信息
	request := &rpc.CertificatePutRequest{
		Name: new.ResID,
		Sni:  []string{new.Sni},
		Cert: new.Certificate,
		Key:  new.PrivateKey,
	}
	err = rpc.NewApiOak().CertificatePut(request)
	if err != nil {
		// 同步删除过数据面旧证书信息时需要回滚，控制面数据根据事务自动回滚
		if tmpDeleteRes != false {
			request := &rpc.CertificatePutRequest{
				Name: existSniCertificate.ResID,
				Sni:  []string{existSniCertificate.Sni},
				Cert: existSniCertificate.Certificate,
				Key:  existSniCertificate.PrivateKey,
			}
			err = rpc.NewApiOak().CertificatePut(request)
			if err != nil {
				packages.Log.Error("rollback old data side certificate error")
				return err
			}
		}
		return err
	}

	return nil
}

//CertificateAdd
func (s *CertificateService) CertificateAdd(request *validators.CertificateAddUpdate) error {
	certificateInfo, err := utils.DiscernCertificate(&request.Certificate)
	if err != nil {
		return err
	}
	err = packages.GetDb().Transaction(func(tx *gorm.DB) error {

		certificates := &models.Certificates{
			Certificate: request.Certificate,
			PrivateKey:  request.PrivateKey,
			ExpiredAt:   certificateInfo.NotAfter,
			Enable:      request.Enable,
			Sni:         request.Sni,
		}

		resID, err := (&models.Certificates{}).CertificatesAdd(tx, certificates)

		if err != nil {
			return err
		}
		certificates.ResID = resID
		// 当前证书设置为启用状态
		if request.Enable == utils.EnableOn {
			err = syncDataSideCertificate(tx, certificates, "")

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// CertificateUpdate
func (s *CertificateService) CertificateUpdate(resID string, request *validators.CertificateAddUpdate) error {

	certificates, err := (&models.Certificates{}).CertificateInfoById(resID)

	if err != nil {
		return errors.New(enums.CodeMessages(enums.CertificateNull))
	}

	discernCertificateInfo, err := utils.DiscernCertificate(&request.Certificate)
	if err != nil {
		return err
	}

	err = packages.GetDb().Transaction(func(tx *gorm.DB) error {

		certificates.Certificate = request.Certificate
		certificates.PrivateKey = request.PrivateKey
		certificates.ExpiredAt = discernCertificateInfo.NotAfter
		certificates.Enable = request.Enable
		certificates.Sni = request.Sni

		err = (&models.Certificates{}).CertificatesUpdate(tx, resID, &certificates)

		if err != nil {
			return err
		}

		if request.Enable == utils.EnableOn {
			err = syncDataSideCertificate(tx, &certificates, certificates.ResID)

			if err != nil {
				return err
			}
		}
		return nil
	})

	return nil
}

type CertificateInfo struct {
	ResID       string `json:"resId"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key"`
	Enable      int    `json:"enable"`
}

// CertificateInfo
func (s *CertificateService) CertificateInfo(id string) (CertificateInfo, error) {

	certificatesModel := models.Certificates{}

	certificateInfo, err := certificatesModel.CertificateInfoById(id)

	if err != nil {
		return CertificateInfo{}, errors.New(enums.CodeMessages(enums.CertificateNull))
	}

	return CertificateInfo{
		ResID:       certificateInfo.ResID,
		Certificate: certificateInfo.Certificate,
		PrivateKey:  certificateInfo.PrivateKey,
		Enable:      certificateInfo.Enable,
	}, nil

}

type CertificateItem struct {
	ResID     string `json:"resId"`
	Sni       string `json:"sni"`
	ExpiredAt int64  `json:"expired_at"`
	Enable    int    `json:"enable"`
}

// CertificateListPage
func (s *CertificateService) CertificateListPage(param *validators.CertificateList) ([]CertificateItem, int, error) {
	certificatesModel := models.Certificates{}
	certificateList, total, err := certificatesModel.CertificateListPage(param)

	if err != nil {
		return []CertificateItem{}, 0, err
	}

	if len(certificateList) == 0 {
		return []CertificateItem{}, 0, nil
	}
	list := []CertificateItem{}

	for _, v := range certificateList {
		list = append(list, CertificateItem{
			ResID:     v.ResID,
			Sni:       v.Sni,
			ExpiredAt: v.ExpiredAt.Unix(),
			Enable:    v.Enable,
		})
	}

	return list, total, nil
}

// CertificateDelete
func (s *CertificateService) CertificateDelete(resID string) error {

	_, err := (&models.Certificates{}).CertificateInfoById(resID)

	if err != nil {
		return errors.New(enums.CodeMessages(enums.CertificateNull))
	}

	err = packages.GetDb().Transaction(func(tx *gorm.DB) error {

		err := (&models.Certificates{}).CertificateDelete(tx, resID)

		if err != nil {
			return err
		}

		err = rpc.NewApiOak().CertificateDelete(resID)

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

// CertificateSwitchEnable
func (s *CertificateService) CertificateSwitchEnable(resID string, enable int) error {

	certificates, err := (&models.Certificates{}).CertificateInfoById(resID)

	if err != nil {
		return errors.New(enums.CodeMessages(enums.CertificateNull))
	}

	err = packages.GetDb().Transaction(func(tx *gorm.DB) error {
		err = (&models.Certificates{}).CertificateSwitchEnable(tx, resID, enable)

		if err != nil {
			return err
		}

		if enable == utils.EnableOn {
			err = syncDataSideCertificate(tx, &certificates, certificates.ResID)

			if err != nil {
				return err
			}
		} else {
			err = rpc.NewApiOak().CertificateDelete(resID)

			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
