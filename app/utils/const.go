package utils

const (
	IdTypeUser          = "us"
	IdTypeUserToken     = "ut"
	IdTypeService       = "sv"
	IdTypeServiceDomain = "sd"
	IdTypeServiceNode   = "sn"
	IdTypeRouter        = "rt"
	IdTypePlugin        = "pl"
	IdTypeRoutePlugin   = "rp"
	IdTypeCertificate   = "ce"
	IdTypeClusterNode   = "cn"
	IdTypeUpstream      = "up"
	IdTypeUpstreamNode  = "un"

	EtcdKeyTypeService     = "service"
	EtcdKeyTypeRouter      = "router"
	EtcdKeyTypePlugin      = "plugin"
	EtcdKeyTypeCertificate = "certificate"

	IdLength           = 15
	IdGenerateMaxTimes = 5

	IPV4 = "ipv4"
	IPV6 = "ipv6"

	IPTypeV4 = 1
	IPTypeV6 = 2

	LocalEn = "en"
	LocalZh = "zh"

	Page     = 1
	PageSize = 10

	MaxPageSize = 100

	EnableOn  = 1
	EnableOff = 2

	EtcdTimeOut = 3

	ReleaseY = 1 // 是否发布——是
	ReleaseN = 2 // 是否发布——否

	ReleaseStatusU = 1 // 发布状态——未发布
	ReleaseStatusT = 2 // 发布状态——待发布
	ReleaseStatusY = 3 // 发布状态——已发布

	ReleaseTypePush   = "push"   // 发布类型——发布（新增/修改）
	ReleaseTypeDelete = "delete" // 发布类型——删除

	// ===================================== upstream =====================================

	LoadBalanceRoundRobin = 1 // 加权轮询 (Round Robin)
	LoadBalanceCHash      = 2 // 一致性Hash（CHash）

	LoadBalanceNameRoundRobin = "加权轮询 (Round Robin)"
	LoadBalanceNameCHash      = "一致性Hash（CHash）"

	ConfigBalanceNameRoundRobin = "ROUNDROBIN"
	ConfigBalanceNameCHash      = "CHASH"

	ProtocolHTTP         = 1
	ProtocolHTTPS        = 2
	ProtocolHTTPAndHTTPS = 3

	// ===================================== upstream node =====================================

	DefaultNodePort = 80

	HealthNodeWeight = 1 // 节点默认权重

	HealthY = 1 // 健康状态——健康
	HealthN = 2 // 健康状态——异常

	HealthNameY = "健康" // 健康状态——健康
	HealthNameN = "异常" // 健康状态——异常

	ConfigHealthY = "HEALTH"   // 健康
	ConfigHealthN = "UNHEALTH" // 异常

	HealthCheckOn  = 1 // 健康检查——开
	HealthCheckOff = 2 // 健康检查——关

	ConfigHealthCheckOn  = true  // 健康检查——开
	ConfigHealthCheckOff = false // 健康检查——关

	// ===================================== route =====================================

	DefaultRouterPath = "/*"

	RequestMethodALL = "ALL"

	RequestMethodGET     = "GET"
	RequestMethodPUT     = "PUT"
	RequestMethodPOST    = "POST"
	RequestMethodPATH    = "PATH"
	RequestMethodDELETE  = "DELETE"
	RequestMethodOPTIONS = "OPTIONS"

	// ===================================== plugin =====================================

	PluginTypeIdAuth        = 1
	PluginTypeIdLimit       = 2
	PluginTypeIdSafety      = 3
	PluginTypeIdFlowControl = 4
	PluginTypeIdOther       = 99

	PluginTypeNameAuth        = "鉴权"
	PluginTypeNameLimit       = "限流"
	PluginTypeNameSafety      = "安全"
	PluginTypeNameFlowControl = "流量控制"
	PluginTypeNameOther       = "其他"

	PluginIdCors       = "pl-dIhZpgqcCHQzNgT"
	PluginIdMock       = "pl-5xO9hzfcHJtpcQT"
	PluginIdKeyAuth    = "pl-xZjvnLQfq2i5GTS"
	PluginIdJwtAuth    = "pl-0FnmajmiO7C8PtX"
	PluginIdLimitReq   = "pl-m5BzSXbCQfGzoQi"
	PluginIdLimitConn  = "pl-rLYsoeNVfPUMUAA"
	PluginIdLimitCount = "pl-XZxaqOgRZsBKpoE"

	PluginKeyCors       = "cors"
	PluginKeyMock       = "mock"
	PluginKeyKeyAuth    = "key-auth"
	PluginKeyJwtAuth    = "jwt-auth"
	PluginKeyLimitReq   = "limit-req"
	PluginKeyLimitConn  = "limit-conn"
	PluginKeyLimitCount = "limit-count"

	PluginIconCors       = "icon-cors"
	PluginIconMock       = "icon-mock"
	PluginIconKeyAuth    = "icon-key-auth"
	PluginIconJwtAuth    = "icon-jwt-auth"
	PluginIconLimitReq   = "icon-limit-req"
	PluginIconLimitConn  = "icon-limit-conn"
	PluginIconLimitCount = "icon-limit-count"

	// ===================================== cluster node =====================================

	EtcdKeyWatchClusterNode = "/apioak/etcd-key/watch/cluster/node/add"

	ClusterNodeStatusHealth    = 1
	ClusterNodeStatusUnhealthy = 2
)
