package providers

import (
	"strings"

	"go.uber.org/zap"

	"github.com/kentalee/errors"

	"github.com/kentalee/ddns/common/defines"
	"github.com/kentalee/ddns/internal/common"
)

var (
	EMInvalidDomainProvider      = errors.New("9caba9e400040001", "无效的域名提供商")
	EMFailedToInitSession        = errors.New("9caba9e400040002", "未能成功创建会话")
	EMDomainUnavailable          = errors.New("9caba9e400040003", "域名不可用")
	ErrFailedToListDomainRecords = errors.New("9caba9e400040004", "FailedToListDomainRecords")
)

type Provider interface {
	InitSession(domain string, config common.STDomain) error
	Update(ipType, address string) error
}

func NewProvider(name string) (Provider, error) {
	switch name {
	case AliyunName:
		return new(Aliyun), nil
	default:
		return nil, errors.Note(EMInvalidDomainProvider, zap.String("provider", name))
	}
}

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
