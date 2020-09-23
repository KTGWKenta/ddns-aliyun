package config

import (
	"gitlab.com/MGEs/Base/contexts"
	"gitlab.com/MGEs/Base/workflow"
	cued "gitlab.com/MGEs/CUEd"
)

type STDomain struct {
	Provider string              `cued:"provider"`
	AuthArgs map[string]string   `cued:"authArgs"`
	Prefixes []map[string]string `cued:"prefixes"`
}

type STLookups struct {
	V4Addr string `cued:"v4Addr"`
	V4Path string `cued:"v4Path"`
	V6Addr string `cued:"v6Addr"`
	V6Path string `cued:"v6Path"`
}

type STConfig struct {
	Domains map[string]STDomain
	Lookups STLookups
}

var Config = STConfig{
	Domains: map[string]STDomain{},
	Lookups: STLookups{},
}

func init() {
	var err error
	var configCtx contexts.Context
	if configCtx, err = contexts.New(&Config.Domains); err != nil {
		workflow.Throw(err, workflow.TlPanic)
	} else if err = cued.RegisterField(cued.Field{
		Path:    []string{"domains"},
		Context: configCtx,
	}); err != nil {
		workflow.Throw(err, workflow.TlPanic)
	}

	if configCtx, err = contexts.New(&Config.Lookups); err != nil {
		workflow.Throw(err, workflow.TlPanic)
	} else if err = cued.RegisterField(cued.Field{
		Path:    []string{"lookups"},
		Context: configCtx,
	}); err != nil {
		workflow.Throw(err, workflow.TlPanic)
	}
}
