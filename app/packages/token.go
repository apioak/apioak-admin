package packages

type tokenConfig struct {
	TokenIssuer string
	TokenSecret string
	TokenExpire uint32
}

var Token tokenConfig

func SetToken(issuer string, secret string, expire uint32) {
	token := tokenConfig{
		TokenIssuer: issuer,
		TokenSecret: secret,
		TokenExpire: expire,
	}
	Token = token
}
