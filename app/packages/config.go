package packages

type configApiOak struct {
	Ip     string
	Port   int
	Domain string
	Secret string
}

var ConfigApiOak configApiOak

func SetConfigApiOak(ip string, port int, domain string, secret string) {
	apiOak := configApiOak{
		Ip:     ip,
		Port:   port,
		Domain: domain,
		Secret: secret,
	}

	ConfigApiOak = apiOak
}
