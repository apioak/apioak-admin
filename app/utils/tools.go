package utils

import (
	"apioak-admin/app/enums"
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
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
	case IdTypeService:
		id = IdTypeService + "-" + randomId
	case IdTypeServiceDomain:
		id = IdTypeServiceDomain + "-" + randomId
	case IdTypeServiceNode:
		id = IdTypeServiceNode + "-" + randomId
	case IdTypeRoute:
		id = IdTypeRoute + "-" + randomId
	case IdTypePlugin:
		id = IdTypePlugin + "-" + randomId
	case IdTypeRoutePlugin:
		id = IdTypeRoutePlugin + "-" + randomId
	case IdTypeCertificate:
		id = IdTypeCertificate + "-" + randomId
	case IdTypeClusterNode:
		id = IdTypeClusterNode + "-" + randomId
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
		{Id: LoadBalanceIPHash, Name: LoadBalanceNameIPHash},
	}

	return loadBalanceList
}

func PluginAllTypes() []enumInfo {
	pluginTypeList := []enumInfo{
		{Id: PluginTypeIdAuth, Name: PluginTypeNameAuth},
		{Id: PluginTypeIdLimit, Name: PluginTypeNameLimit},
	}

	return pluginTypeList
}

func AllRequestMethod() []string {
	return []string{
		RequestMethodALL,
		RequestMethodGET,
		RequestMethodPOST,
		RequestMethodPUT,
		RequestMethodDELETE,
		RequestMethodOPTIONS,
	}
}
