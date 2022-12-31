package cores

import (
	"apioak-admin/app/packages"
	"errors"
)

func InitToken(conf *ConfigGlobal) error {
	issuer := conf.Token.TokenIssuer
	secret := conf.Token.TokenSecret
	expire := conf.Token.TokenExpire

	if len(issuer) <= 0 {
		return errors.New("The token issuer configuration cannot be empty")
	}
	if len(secret) <= 0 {
		return errors.New("The token secret configuration cannot be empty")
	}
	if expire < 1 {
		return errors.New("The minimum token secret can only be 1")
	}

	if expire > 120 {
		expire = 120
	}

	packages.SetToken(issuer, secret, expire)

	return nil
}
