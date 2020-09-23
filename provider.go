package main

import (
	"github.com/KTGWKenta/ddns-aliyun/config"
	"gitlab.com/MGEs/Base/workflow"

	"github.com/KTGWKenta/ddns-aliyun/providers"
)

var EMInvalidDomainProvider = workflow.NewMask(
	"EMInvalidDomainProvider",
	"无效的域名提供商`{{provider}}`",
)

type provider interface {
	InitSession(domain string, config config.STDomain) error
	Update(ipType, address string) error
}

func newProvider(name string) (provider, error) {
	switch name {
	case providers.AliyunName:
		return new(providers.Aliyun), nil
	default:
		return nil, workflow.NewException(
			EMInvalidDomainProvider, map[string]string{"provider": name}, nil,
		)
	}
}
