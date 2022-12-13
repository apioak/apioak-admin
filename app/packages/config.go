package packages

type configApiOak struct {
	Protocol string
	Ip       string
	Port     int
	Domain   string
	Secret   string
}

var ConfigApiOak configApiOak

func SetConfigApiOak(protocol string, ip string, port int, domain string, secret string) {
	apiOak := configApiOak{
		Protocol: protocol,
		Ip:       ip,
		Port:     port,
		Domain:   domain,
		Secret:   secret,
	}

	ConfigApiOak = apiOak
}
