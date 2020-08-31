package main

import (
	"gitlab.com/MGEs/Base/contexts"
	"gitlab.com/MGEs/Base/workflow"
	cued "gitlab.com/MGEs/CUEd"
)

type STDomain struct {
	Provider string
	AuthArgs map[string]string
}

type STConfig struct {
	Domains []STDomain
}

var Config = STConfig{
	Domains: []STDomain{},
}

func initConfig() {
	var err error
	var configCtx contexts.Context
	if configCtx, err = contexts.New(&Config.Domains); err != nil {
		workflow.Throw(err, workflow.TlPanic)
	} else if err = cued.RegisterField(cued.Field{
		Path:    []string{"Domains"},
		Context: configCtx,
	}); err != nil {
		workflow.Throw(err, workflow.TlPanic)
	}
}
