package validators

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var (
	domainLocalErrorMessages = map[string]map[string]string{
		utils.LocalEn: {
			"required": " is a required field",
			"max":      " must be %d or less",
			"min":      " must be %d or greater",
		},
		utils.LocalZh: {
			"required": "为必填字段",
			"max":      "长度必须小于或等于%d",
			"min":      "长度最小只能为%d",
		},
	}
	domainMin = 1
	domainMax = 50
)

type ServiceDomainAddUpdate struct {
	Domain string `json:"service_domain" zh:"单个域名" en:"Domain name" binding:"required,min=1,max=50"`
}

func CheckServiceDomain(fl validator.FieldLevel) bool {
	serviceDomainInterface := fl.Field().Interface()
	serviceDomains := serviceDomainInterface.([]string)

	serviceDomainValidator := packages.GetCustomizeValidator()
	for _, domain := range serviceDomains {
		domainTrim := strings.TrimSpace(domain)

		serviceDomain := ServiceDomainAddUpdate{
			Domain: domainTrim,
		}

		domainErr := serviceDomainValidator.Struct(serviceDomain)
		if domainErr != nil {
			var (
				structField string
				tag         string
				field       string
				errMsg      string
			)

			for _, e := range domainErr.(validator.ValidationErrors) {
				structField = e.StructField()
				tag = e.Tag()
				field = e.Field()
				break
			}

			switch strings.ToLower(structField) {
			case "domain":
				errMsg = domainValidator(tag, field)
			}
			packages.SetAllCustomizeValidatorErrMsgs("CheckServiceDomain", errMsg)
			return false
		}
	}
	return true
}

func domainValidator(tag string, field string) string {
	var errMsg string

	switch strings.ToLower(tag) {
	case "required":
		errMsg = fmt.Sprintf(field + domainLocalErrorMessages[strings.ToLower(packages.GetValidatorLocale())][strings.ToLower(tag)])
	case "min":
		errMsg = fmt.Sprintf(field+domainLocalErrorMessages[strings.ToLower(packages.GetValidatorLocale())][strings.ToLower(tag)], domainMin)
	case "max":
		errMsg = fmt.Sprintf(field+domainLocalErrorMessages[strings.ToLower(packages.GetValidatorLocale())][strings.ToLower(tag)], domainMax)
	}
	return errMsg
}

func CorrectServiceDomains(serviceDomains []string) {
	tmpDomainsMap := make(map[string]byte, 0)
	tmpServiceDomains := make([]string, 0)

	for _, domain := range serviceDomains {
		domainTrim := strings.TrimSpace(domain)
		if len(domainTrim) <= 0 {
			continue
		}

		_, exist := tmpDomainsMap[domainTrim]
		if exist {
			continue
		}

		tmpDomainsMap[domainTrim] = 0
		tmpServiceDomains = append(tmpServiceDomains, domainTrim)
	}

	serviceDomains = tmpServiceDomains
}
