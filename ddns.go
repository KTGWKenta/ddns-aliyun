package main

import (
	"errors"

	"github.com/KTGWKenta/ddns-aliyun/config"
	"github.com/KTGWKenta/ddns-aliyun/defines"
	"gitlab.com/MGEs/Base/workflow"
)

func payloadDispatcher(ipType string, handler func(ipType, address string) error) func(event workflow.Event) error {
	return func(event workflow.Event) error {
		if address, ok := event.GetPayload().(string); !ok {
			return errors.New("invalid event payload")
		} else {
			return handler(ipType, address)
		}
	}
}

func applyConfigs(event workflow.Event) error {
	var err error
	for domain, cfg := range config.Config.Domains {
		var p provider
		if p, err = newProvider(cfg.Provider); err != nil {
			workflow.Throw(err, workflow.TlError)
			continue
		}
		if err = p.InitSession(domain, cfg); err != nil {
			workflow.Throw(err, workflow.TlError)
			continue
		}
		err = workflow.GlobalEvents().Subscribe(defines.EVTUpdateIPV4, 0, domain, payloadDispatcher(defines.IPTypeV4, p.Update))
		if err != nil {
			workflow.Throw(err, workflow.TlError)
			continue
		}
		err = workflow.GlobalEvents().Subscribe(defines.EVTUpdateIPV6, 0, domain, payloadDispatcher(defines.IPTypeV6, p.Update))
		if err != nil {
			workflow.Throw(err, workflow.TlError)
			continue
		}
	}
	return nil
}

func initDDNS() {
	if err := workflow.GlobalEvents().Subscribe(
		workflow.EVT_Lifecycle_Initialize, workflow.EventPriority_Latest, "ddns", applyConfigs,
	); err != nil {
		workflow.Throw(err, workflow.TlPanic)
	}
}
