package providers

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"go.uber.org/zap"

	"github.com/kentalee/errors"
	"github.com/kentalee/log"

	"github.com/kentalee/ddns-aliyun/common/defines"
	"github.com/kentalee/ddns-aliyun/internal/common"
)

const (
	AliyunName                  = "aliyun"
	AliyunAuthFieldRegion       = "region"
	AliyunAuthFieldAccessKey    = "accessKey"
	AliyunAuthFieldAccessSecret = "accessSecret"
	AliyunPrefixFieldLine       = "line"
	AliyunPrefixFieldIpType     = "ipType"
	AliyunPrefixFieldRecord     = "record"
	AliyunPrefixFieldRecordType = "recordType"
	AliyunAPIMaxPageSize        = 250
)

type AliyunRecord struct {
	ipType     string
	record     string
	recordType string
	line       string
}

type Aliyun struct {
	client  *alidns.Client
	domain  string
	records []*AliyunRecord
}

func (p *Aliyun) InitSession(domain string, config common.STDomain) error {
	var err error
	var client *alidns.Client
	var records []*AliyunRecord
	if client, err = alidns.NewClientWithAccessKey(
		config.AuthArgs[AliyunAuthFieldRegion],
		config.AuthArgs[AliyunAuthFieldAccessKey],
		config.AuthArgs[AliyunAuthFieldAccessSecret],
	); err != nil {
		return errors.Because(EMFailedToInitSession, err,
			zap.String("keyName", "accessKey"),
			zap.String("key", config.AuthArgs[AliyunAuthFieldAccessKey]),
		)
	}
	for i, prefix := range config.Prefixes {
		if record, err := p.parseSubDomains(prefix); err != nil {
			return errors.Note(err, zap.Int("prefix index", i))
		} else {
			records = append(records, record)
		}
	}
	p.records = records
	p.client = client
	p.domain = domain
	return p.checkDomain()
}

func (p *Aliyun) parseSubDomains(val map[string]string) (*AliyunRecord, error) {
	recordObj := AliyunRecord{}
	if ipType, ok := val[AliyunPrefixFieldIpType]; ok {
		if ipType != defines.IPTypeV4 && ipType != defines.IPTypeV6 {
			return nil, errors.NewSysErr("invalid ip type")
		}
		recordObj.ipType = ipType
	} else {
		return nil, errors.NewSysErr("ipType is required")
	}
	if record, ok := val[AliyunPrefixFieldRecord]; ok {
		recordObj.record = record
	} else {
		return nil, errors.NewSysErr("record is required")
	}
	if recordType, ok := val[AliyunPrefixFieldRecordType]; ok {
		recordObj.recordType = recordType
	} else {
		return nil, errors.NewSysErr("recordType is required")
	}
	if line, ok := val[AliyunPrefixFieldLine]; ok {
		recordObj.line = line
	} else {
		recordObj.line = "default"
	}
	return &recordObj, nil
}

func (p *Aliyun) checkDomain() error {
	req := alidns.CreateDescribeDomainInfoRequest()
	req.DomainName = p.domain
	response, err := p.client.DescribeDomainInfo(req)
	if err != nil {
		return errors.Because(EMDomainUnavailable, err, zap.String("domain", p.domain))
	}
	log.With(
		"id", response.DomainId,
		"domain", response.DomainName,
	).Info("domainCheckResp")
	return nil
}

func (p *Aliyun) getRecords() (map[string][]alidns.Record, error) {
	const pageSize = int64(AliyunAPIMaxPageSize)
	lsReq := alidns.CreateDescribeDomainRecordsRequest()
	lsReq.DomainName = p.domain
	lsReq.Status = "Enable"
	lsReq.PageSize = requests.NewInteger(AliyunAPIMaxPageSize)
	var total = pageSize
	var records = map[string][]alidns.Record{}
	for i := int64(1); i*pageSize <= total; i++ {
		lsResp, err := p.client.DescribeDomainRecords(lsReq)
		if err != nil {
			return nil, errors.Note(ErrFailedToListDomainRecords,
				zap.Int64("pn", i),
				zap.Int64("size", pageSize),
				zap.String("domain", p.domain),
			)
		}
		total = lsResp.TotalCount
		for _, record := range lsResp.DomainRecords.Record {
			records[record.RR] = append(records[record.RR], record)
		}
	}
	return records, nil
}

var (
	ErrFailedToCreateRecord = errors.New("9caba9e400040005", "failed to create record")
	ErrFailedToDeleteRecord = errors.New("9caba9e400040006", "failed to delete record")
	ErrFailedToUpdateRecord = errors.New("9caba9e400040007", "failed to update record")
)

func (p *Aliyun) Update(ipType, address string) error {
	log.With("type", ipType, "addr", address).Info("newIP")
	var err error
	var records map[string][]alidns.Record
	if records, err = p.getRecords(); err != nil {
		return err
	}
	for _, prefix := range p.records {
		if prefix.ipType != ipType {
			continue
		}
		if recordGroup, exists := records[prefix.record]; exists {
			var NormalRecords = make([]alidns.Record, 0)
			var LockedRecords = make([]alidns.Record, 0)
			for _, record := range recordGroup {
				if record.Line == prefix.line && record.Type == prefix.recordType {
					if record.Locked {
						LockedRecords = append(LockedRecords, record)
					} else {
						NormalRecords = append(NormalRecords, record)
					}
				}
			}
			if len(NormalRecords)+len(LockedRecords) > 1 || len(LockedRecords) == 1 {
				// delete
				var batchArr []alidns.OperateBatchDomainDomainRecordInfo
				for _, record := range NormalRecords {
					batchArr = append(batchArr, alidns.OperateBatchDomainDomainRecordInfo{
						Type:   "RR_DEL",
						Rr:     prefix.record,
						Value:  record.Value,
						Domain: p.domain,
					})
				}
				req := alidns.CreateOperateBatchDomainRequest()
				req.DomainRecordInfo = &batchArr
				if _, err := p.client.OperateBatchDomain(req); err != nil {
					log.Error(errors.Note(ErrFailedToDeleteRecord, zap.String("domain", p.domain)))
				}
			} else if len(NormalRecords) == 1 {
				if NormalRecords[0].Value != address {
					// update
					req := alidns.CreateUpdateDomainRecordRequest()
					req.RecordId = NormalRecords[0].RecordId
					req.RR = prefix.record
					req.Value = address
					req.Type = prefix.recordType
					if _, err := p.client.UpdateDomainRecord(req); err != nil {
						log.Error(errors.Note(ErrFailedToUpdateRecord, zap.String("domain", p.domain)))
					}
				}
				continue
			}
		}
		// add
		req := alidns.CreateAddDomainRecordRequest()
		req.DomainName = p.domain
		req.RR = prefix.record
		req.Value = address
		req.Type = prefix.recordType
		if _, err := p.client.AddDomainRecord(req); err != nil {
			log.Error(errors.Because(ErrFailedToCreateRecord, err,
				zap.String("domain", p.domain),
				zap.String("address", address),
				zap.String("type", prefix.recordType)))
		}
	}
	return nil
}
