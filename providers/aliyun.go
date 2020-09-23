package providers

import (
	"strconv"

	"github.com/KTGWKenta/ddns-aliyun/config"
	"github.com/KTGWKenta/ddns-aliyun/defines"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/pkg/errors"
	"gitlab.com/MGEs/Base/workflow"
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

func (p *Aliyun) InitSession(domain string, config config.STDomain) error {
	var err error
	var client *alidns.Client
	var records []*AliyunRecord
	if client, err = alidns.NewClientWithAccessKey(
		config.AuthArgs[AliyunAuthFieldRegion],
		config.AuthArgs[AliyunAuthFieldAccessKey],
		config.AuthArgs[AliyunAuthFieldAccessSecret],
	); err != nil {
		return workflow.NewException(EMFailedToInitSession,
			map[string]string{"keyName": "accessKey", "key": config.AuthArgs[AliyunAuthFieldAccessKey]}, err)
	}
	for i, prefix := range config.Prefixes {
		if record, err := p.parseSubDomains(prefix); err != nil {
			return errors.Wrapf(err, "failed to parse prefix #%d", i)
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
			return nil, errors.New("invalid ip type")
		}
		recordObj.ipType = ipType
	} else {
		return nil, errors.New("ipType is required")
	}
	if record, ok := val[AliyunPrefixFieldRecord]; ok {
		recordObj.record = record
	} else {
		return nil, errors.New("record is required")
	}
	if recordType, ok := val[AliyunPrefixFieldRecordType]; ok {
		recordObj.recordType = recordType
	} else {
		return nil, errors.New("recordType is required")
	}
	if line, ok := val[AliyunPrefixFieldLine]; ok {
		recordObj.line = line
	} else {
		recordObj.line = "default"
	}
	return &recordObj, nil
}

var domainCheckResp = workflow.NewMask("domainCheckResp", "Domain: {{domain}}, Id: {{id}}")
var newIPMsg = workflow.NewMask("newIP", "new ip{{type}}, {{addr}}")

func (p *Aliyun) checkDomain() error {
	req := alidns.CreateDescribeDomainInfoRequest()
	req.DomainName = p.domain
	response, err := p.client.DescribeDomainInfo(req)
	if err != nil {
		return workflow.NewException(EMDomainUnavailable,
			map[string]string{"domain": p.domain}, err)
	}
	workflow.Throw(workflow.NewSimpleThrowable(domainCheckResp, map[string]string{
		"id":     response.DomainId,
		"domain": response.DomainName,
	}), workflow.TlNotify)
	return nil
}

var EM_FailedToListDomainRecords = workflow.NewMask(
	"FailedToListDomainRecords",
	"failed to list domain records: {{domain}}, type:{{type}}, pn: {{pn}}, size:{{size}}",
)

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
			return nil, workflow.NewSimpleThrowable(EM_FailedToListDomainRecords, map[string]string{
				"pn":     strconv.FormatInt(i, 10),
				"size":   strconv.FormatInt(pageSize, 10),
				"domain": p.domain,
			})
		}
		total = lsResp.TotalCount
		for _, record := range lsResp.DomainRecords.Record {
			records[record.RR] = append(records[record.RR], record)
		}
	}
	return records, nil
}

func (p *Aliyun) Update(ipType, address string) error {
	workflow.Throw(workflow.NewSimpleThrowable(newIPMsg, map[string]string{
		"type": ipType,
		"addr": address,
	}), workflow.TlNotify)

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
					workflow.Throw(errors.Errorf("failed to delete record for `%s`:%v", p.domain, err), workflow.TlError)
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
						workflow.Throw(errors.Errorf("failed to update record for `%s`:%v", p.domain, err), workflow.TlError)
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
			workflow.Throw(errors.Errorf("failed to add record for `%s`:%v", p.domain, err), workflow.TlError)
		}
	}
	return nil
}
