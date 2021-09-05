package utils

const (
	IdTypeUser          = "u"
	IdTypeService       = "svc"
	IdTypeServiceDomain = "sdm"
	IdTypeServiceNode   = "snd"
	IdTypeRoute         = "rt"
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

	LoadBalanceRoundRobin = 1 // 轮询
	LoadBalanceIPHash     = 2 // ip_hash

	LoadBalanceNameRoundRobin = "加权轮询 (Weighted Round Robin)" // 轮询
	LoadBalanceNameIPHash     = "ip_hash"                     // ip_hash

	ProtocolHTTP         = 1
	ProtocolHTTPS        = 2
	ProtocolHTTPAndHTTPS = 3
)
