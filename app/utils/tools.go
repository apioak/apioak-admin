package utils

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
	"time"
)

func createRandomString(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}

func IdGenerate(idType string) (string, error) {
	randomId := createRandomString(IdLength)

	var id string
	switch strings.ToLower(idType) {
	case IdTypeUser:
		id = IdTypeUser + "-" + randomId
	case IdTypeUserToken:
		id = IdTypeUserToken + "-" + randomId
	case IdTypeService:
		id = IdTypeService + "-" + randomId
	case IdTypeServiceDomain:
		id = IdTypeServiceDomain + "-" + randomId
	case IdTypeServiceNode:
		id = IdTypeServiceNode + "-" + randomId
	case IdTypeRouter:
		id = IdTypeRouter + "-" + randomId
	case IdTypePlugin:
		id = IdTypePlugin + "-" + randomId
	case IdTypeRoutePlugin:
		id = IdTypeRoutePlugin + "-" + randomId
	case IdTypeCertificate:
		id = IdTypeCertificate + "-" + randomId
	case IdTypeClusterNode:
		id = IdTypeClusterNode + "-" + randomId
	case IdTypeUpstream:
		id = IdTypeUpstream + "-" + randomId
	case IdTypeUpstreamNode:
		id = IdTypeUpstreamNode + "-" + randomId
	default:
		return "", fmt.Errorf("id type error")
	}

	return id, nil
}

func DiscernIP(s string) (string, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return "", fmt.Errorf("(%s) is illegal ip", s)
	}

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return IPV4, nil
		case ':':
			return IPV6, nil
		}
	}
	return "", nil
}

type CertificateInfo struct {
	CommonName string
	NotBefore  time.Time
	NotAfter   time.Time
}

func DiscernCertificate(certificate *string) (CertificateInfo, error) {
	certificateInfo := CertificateInfo{}
	pemBlock, _ := pem.Decode([]byte(*certificate))
	if pemBlock == nil {
		return certificateInfo, errors.New(enums.CodeMessages(enums.CertificateFormatError))
	}

	parseCert, parseCertErr := x509.ParseCertificate(pemBlock.Bytes)
	if parseCertErr != nil {
		return certificateInfo, errors.New(enums.CodeMessages(enums.CertificateParseError))
	}

	certificateInfo.CommonName = parseCert.Subject.CommonName
	certificateInfo.NotBefore = parseCert.NotBefore
	certificateInfo.NotAfter = parseCert.NotAfter

	return certificateInfo, nil
}

type enumInfo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func LoadBalanceList() []enumInfo {
	loadBalanceList := []enumInfo{
		{Id: LoadBalanceRoundRobin, Name: LoadBalanceNameRoundRobin},
		{Id: LoadBalanceCHash, Name: LoadBalanceNameCHash},
	}

	return loadBalanceList
}

func ConfigBalanceList() []enumInfo {
	configBalanceList := []enumInfo{
		{Id: LoadBalanceRoundRobin, Name: ConfigBalanceNameRoundRobin},
		{Id: LoadBalanceCHash, Name: ConfigBalanceNameCHash},
	}

	return configBalanceList
}

func ConfigUpstreamNodeHealthList() []enumInfo {
	configHealthList := []enumInfo{
		{Id: HealthY, Name: ConfigHealthY},
		{Id: HealthN, Name: ConfigHealthN},
	}

	return configHealthList
}

func PluginAllTypes() []enumInfo {
	pluginTypeList := []enumInfo{
		{Id: PluginTypeIdAuth, Name: PluginTypeNameAuth},
		{Id: PluginTypeIdLimit, Name: PluginTypeNameLimit},
		{Id: PluginTypeIdSafety, Name: PluginTypeNameSafety},
		{Id: PluginTypeIdFlowControl, Name: PluginTypeNameFlowControl},
		{Id: PluginTypeIdOther, Name: PluginTypeNameOther},
	}

	return pluginTypeList
}

func PluginAllKeys() []string {
	pluginKeysList := []string{
		PluginKeyCors,
		PluginKeyMock,
		PluginKeyKeyAuth,
		PluginKeyJwtAuth,
		PluginKeyLimitReq,
		PluginKeyLimitConn,
		PluginKeyLimitCount,
	}

	return pluginKeysList
}

type ConfigPluginData struct {
	ResID       string
	PluginKey   string
	Icon        string
	Type        int
	Description string
}

func AllConfigPluginData() []ConfigPluginData {

	allConfigPluginData := []ConfigPluginData{
		{
			ResID: PluginIdCors, PluginKey: PluginKeyCors,
			Icon: PluginIconCors, Type: PluginTypeIdSafety,
		},
		{
			ResID: PluginIdMock, PluginKey: PluginKeyMock,
			Icon: PluginIconMock, Type: PluginTypeIdOther,
		},
		{
			ResID: PluginIdKeyAuth, PluginKey: PluginKeyKeyAuth,
			Icon: PluginIconKeyAuth, Type: PluginTypeIdAuth,
		},
		{
			ResID: PluginIdJwtAuth, PluginKey: PluginKeyJwtAuth,
			Icon: PluginIconJwtAuth, Type: PluginTypeIdAuth,
		},
		{
			ResID: PluginIdLimitReq, PluginKey: PluginKeyLimitReq,
			Icon: PluginIconLimitReq, Type: PluginTypeIdLimit,
		},
		{
			ResID: PluginIdLimitConn, PluginKey: PluginKeyLimitConn,
			Icon: PluginIconLimitConn, Type: PluginTypeIdLimit,
		},
		{
			ResID: PluginIdLimitCount, PluginKey: PluginKeyLimitCount,
			Icon: PluginIconLimitCount, Type: PluginTypeIdLimit,
		},
	}

	return allConfigPluginData
}

func AllRequestMethod() []string {
	return []string{
		RequestMethodALL,
		RequestMethodGET,
		RequestMethodPOST,
		RequestMethodPUT,
		RequestMethodPATH,
		RequestMethodDELETE,
		RequestMethodOPTIONS,
	}
}

func ConfigAllRequestMethod() []string {
	return []string{
		RequestMethodGET,
		RequestMethodPOST,
		RequestMethodPUT,
		RequestMethodPATH,
		RequestMethodDELETE,
		RequestMethodOPTIONS,
	}
}

func Md5(src string) string {
	m := md5.New()
	m.Write([]byte(src))
	srcMd5 := hex.EncodeToString(m.Sum(nil))

	return srcMd5
}

type ExpireToken struct {
	Expire int64
	Token  string
}

type TokenClaims struct {
	Encryption string `json:"encryption"`
	Timestamp  string `json:"timestamp"`
	Secret     string `json:"secret"`
	Issuer     string `json:"issuer"`
}

func GenToken(encryption string) (string, error) {
	tokenClaims := TokenClaims{
		Encryption: encryption,
		Timestamp:  Md5(strconv.FormatInt(time.Now().UnixNano(), 10)),
		Secret:     packages.Token.TokenSecret,
		Issuer:     packages.Token.TokenIssuer,
	}

	var (
		jsonValue []byte
		err       error
	)
	if jsonValue, err = json.Marshal(tokenClaims); err != nil {
		return "", err
	}

	token := strings.TrimRight(base64.URLEncoding.EncodeToString(jsonValue), "=")

	return token, nil
}

func ParseToken(tokenString string) (string, error) {
	if l := len(tokenString) % 4; l > 0 {
		tokenString += strings.Repeat("=", 4-l)
	}

	tokenStructStr, tokenStructStrErr := base64.URLEncoding.DecodeString(tokenString)
	if tokenStructStrErr != nil {
		return "", tokenStructStrErr
	}

	tokenClaims := TokenClaims{}
	unmarshalErr := json.Unmarshal(tokenStructStr, &tokenClaims)
	if unmarshalErr != nil {
		return "", unmarshalErr
	}

	if tokenClaims.Issuer != packages.Token.TokenIssuer || tokenClaims.Secret != packages.Token.TokenSecret {
		return "", errors.New("token parsing failed")
	}

	return tokenClaims.Encryption, nil
}

func IPNameToType(ipName string) (int, error) {
	iPNameToTypeMap := map[string]int{
		IPV4: IPTypeV4,
		IPV6: IPTypeV6,
	}

	ipType, ipTypeExist := iPNameToTypeMap[ipName]
	if ipTypeExist == false {
		return -1, errors.New("IP type does not exist")
	}

	return ipType, nil
}

func IpTypeNameList() (list []enumInfo) {
	list = []enumInfo{
		{Id: IPTypeV4, Name: IPV4},
		{Id: IPTypeV6, Name: IPV6},
	}

	return
}

func IpIdNameMap() (ipIdNameMap map[int]string) {
	ipTypeNameList := IpTypeNameList()

	ipIdNameMap = make(map[int]string)
	for _, IpTypeNameDetail := range ipTypeNameList {
		ipIdNameMap[IpTypeNameDetail.Id] = IpTypeNameDetail.Name
	}

	return
}

func IpNameIdMap() (nameIdMap map[string]int) {
	ipTypeNameList := IpTypeNameList()

	nameIdMap = make(map[string]int)
	for _, IpTypeNameDetail := range ipTypeNameList {
		nameIdMap[IpTypeNameDetail.Name] = IpTypeNameDetail.Id
	}

	return
}

func HealthTypeNameList() (list []enumInfo) {
	list = []enumInfo{
		{Id: HealthY, Name: HealthNameY},
		{Id: HealthN, Name: HealthNameN},
	}

	return
}

func HealthTypeNameMap() (healthNameMap map[int]string) {
	healthNameList := HealthTypeNameList()

	healthNameMap = make(map[int]string)
	for _, healthNameDetail := range healthNameList {
		healthNameMap[healthNameDetail.Id] = healthNameDetail.Name
	}

	return
}


func InterceptSni(domains []string) ([]string, error) {
	domainSniInfos := make([]string, 0)
	if len(domains) == 0 {
		return domainSniInfos, nil
	}

	tmpDomainSniMap := make(map[string]byte, 0)
	for _, domain := range domains {
		disassembleDomains := strings.Split(domain, ".")
		if len(disassembleDomains) < 2 {
			return domainSniInfos, errors.New(enums.CodeMessages(enums.ServiceDomainFormatError))
		}

		disassembleDomains[0] = "*"
		domainSniInfo := strings.Join(disassembleDomains, ".")

		_, exit := tmpDomainSniMap[domainSniInfo]
		if exit {
			continue
		}

		tmpDomainSniMap[domainSniInfo] = 0
		domainSniInfos = append(domainSniInfos, domainSniInfo)
	}

	return domainSniInfos, nil
}

func EtcdKey(keyType string, id string) string {
	key := ""
	if (len(keyType) == 0) || (len(id) == 0) {
		return key
	}

	prefix := "/apioak/"
	switch strings.ToLower(keyType) {
	case EtcdKeyTypeService:
		key = prefix + EtcdKeyTypeService + "/" + id
	case EtcdKeyTypeRouter:
		key = prefix + EtcdKeyTypeRouter + "/" + id
	case EtcdKeyTypePlugin:
		key = prefix + EtcdKeyTypePlugin + "/" + id
	case EtcdKeyTypeCertificate:
		key = prefix + EtcdKeyTypeCertificate + "/" + id
	}

	return key
}
