package providers

import (
	"strings"

	"github.com/KTGWKenta/ddns-aliyun/defines"
	"gitlab.com/MGEs/Base/workflow"
)

var EMFailedToInitSession = workflow.NewMask(
	"EMFailedToInitSession",
	"未能成功创建会话: {{keyName}}:{{key}}",
)

var EMDomainUnavailable = workflow.NewMask(
	"EMFailedToInitSession",
	"域名`{{domain}}`不可用",
)

func ipTypeToNSType(ipType string) string {
	switch strings.ToLower(ipType) {
	case defines.IPTypeV4:
		return "A"
	case defines.IPTypeV6:
		return "AAAA"
	default:
		return "TXT"
	}
}