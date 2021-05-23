package internal

import (
	"github.com/kentalee/log"

	"github.com/kentalee/ddns/common/config"
	"github.com/kentalee/ddns/common/events"
	"github.com/kentalee/ddns/internal/common"
	"github.com/kentalee/ddns/internal/providers"
)

var Config = common.STConfig{}

func applyConfigs() error {
	var err error
	cfgVal := config.Lookup()
	if err = cfgVal.Err(); err != nil {
		log.Fatal(err)
	}
	if err = config.Decode(cfgVal, &Config); err != nil {
		log.Fatal(err)
	}
	for domain, cfg := range Config.Domains {
		var p providers.Provider
		if p, err = providers.NewProvider(cfg.Provider); err != nil {
			return err
		}
		if err = p.InitSession(domain, cfg); err != nil {
			return err
		}
		if err = events.Global().Register(&IPV4UpdateHandler{provider: p}); err != nil {
			return err
		}
		if err = events.Global().Register(&IPV4UpdateHandler{provider: p}); err != nil {
			return err
		}
	}
	return nil
}
