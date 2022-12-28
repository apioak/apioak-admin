package validators

type CertificateAddUpdate struct {
	Sni         string `json:"sni" zh:"域名" en:"Domain name" binding:"required"`
	Certificate string `json:"certificate" zh:"证书内容" en:"Certificate content" binding:"required"`
	PrivateKey  string `json:"private_key" zh:"私钥内容" en:"Private key content" binding:"required"`
	Enable      int    `json:"enable" zh:"证书开关" en:"Certificate enable" binding:"required,oneof=1 2"`
}

type CertificateList struct {
	Enable int    `form:"enable" json:"enable" zh:"证书开关" en:"Certificate enable" binding:"omitempty,oneof=1 2"`
	Search string `form:"search" json:"search" zh:"搜索内容" en:"Search content" binding:"omitempty"`
	BaseListPage
}

type CertificateSwitchEnable struct {
	Enable int `form:"enable" json:"enable" zh:"证书开关" en:"Certificate enable" binding:"required,oneof=1 2"`
}
