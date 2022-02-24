package validators

type CertificateAddUpdate struct {
	Certificate string `json:"certificate" zh:"证书内容" en:"Certificate content" binding:"required"`
	PrivateKey  string `json:"private_key" zh:"私钥内容" en:"Private key content" binding:"required"`
	IsEnable    int    `json:"is_enable" zh:"证书开关" en:"Certificate enable" binding:"required,oneof=1 2"`
	IsRelease   int    `json:"is_release" zh:"发布开关" en:"Release status enable" binding:"omitempty,oneof=1 2"`
}

type CertificateList struct {
	IsEnable      int    `form:"is_enable" json:"is_enable" zh:"证书开关" en:"Certificate enable" binding:"omitempty,oneof=1 2"`
	ReleaseStatus int    `form:"release_status" json:"release_status" zh:"发布状态" en:"Release status" binding:"omitempty,oneof=1 2 3"`
	Search        string `form:"search" json:"search" zh:"搜索内容" en:"Search content" binding:"omitempty"`
	BaseListPage
}

type CertificateSwitchEnable struct {
	IsEnable int `form:"is_enable" json:"is_enable" zh:"证书开关" en:"Certificate enable" binding:"required,oneof=1 2"`
}
