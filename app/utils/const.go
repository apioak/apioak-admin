package utils

const (
	IdTypeUser          = "u"
	IdTypeService       = "svc"
	IdTypeServiceDomain = "sdm"
	IdTypeServiceNode   = "snd"
	IdTypeRoute         = "rt"
	IdTypePlugin        = "plu"
	IdTypeRoutePlugin   = "rpu"
	IdTypeCertificate   = "cer"
	IdTypeClusterNode   = "cnd"

	IdLength = 15

	IPV4 = "ipv4"
	IPV6 = "ipv6"

	LocalEn = "en"
	LocalZh = "zh"

	Page     = 1
	PageSize = 10

	MaxPageSize = 100

	EnableOn  = 1
	EnableOff = 2

	// ===================================== service =====================================

	LoadBalanceRoundRobin = 1 // 轮询
	LoadBalanceIPHash     = 2 // ip_hash

	LoadBalanceNameRoundRobin = "加权轮询 (Weighted Round Robin)"
	LoadBalanceNameIPHash     = "ip_hash"

	ProtocolHTTP         = 1
	ProtocolHTTPS        = 2
	ProtocolHTTPAndHTTPS = 3

	// ===================================== route =====================================

	DefaultRoutePath = "/*"

	RequestMethodALL     = "ALL"
	RequestMethodGET     = "GET"
	RequestMethodPOST    = "POST"
	RequestMethodPUT     = "PUT"
	RequestMethodDELETE  = "DELETE"
	RequestMethodOPTIONS = "OPTIONS"
)
