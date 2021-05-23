package internal

import (
	"fmt"

	"github.com/kentalee/eventbus"
	"github.com/kentalee/log"

	"github.com/kentalee/ddns/common/defines"
	"github.com/kentalee/ddns/common/events"
	"github.com/kentalee/ddns/internal/providers"
)

type ipResponse struct {
	Ip string `json:"ip"`
}

type IPV4UpdatePoster struct {
	ip string
}

func (p *IPV4UpdatePoster) Event() eventbus.Key {
	return defines.EVTUpdateIPV4
}

func (p *IPV4UpdatePoster) IP() string {
	return p.ip
}

type IPV6UpdatePoster struct {
	ip string
}

func (p *IPV6UpdatePoster) Event() eventbus.Key {
	return defines.EVTUpdateIPV6
}

func (p *IPV6UpdatePoster) IP() string {
	return p.ip
}

type IPV4UpdateHandler struct {
	provider providers.Provider
}

func (I IPV4UpdateHandler) Events() map[eventbus.Key]eventbus.Priority {
	return map[eventbus.Key]eventbus.Priority{
		defines.EVTUpdateIPV4: 0,
	}
}

func (I IPV4UpdateHandler) OnEvent(key eventbus.Key, publisher eventbus.Publisher) error {
	if key == defines.EVTUpdateIPV4 {
		if poster, ok := publisher.(*IPV4UpdatePoster); ok {
			return I.provider.Update(defines.IPTypeV4, poster.IP())
		} else {
			return fmt.Errorf("invalid ip poster")
		}
	}
	return nil
}

type IPV6UpdateHandler struct {
	provider providers.Provider
}

func (I IPV6UpdateHandler) Events() map[eventbus.Key]eventbus.Priority {
	return map[eventbus.Key]eventbus.Priority{
		defines.EVTUpdateIPV6: 0,
	}
}

func (I IPV6UpdateHandler) OnEvent(key eventbus.Key, publisher eventbus.Publisher) error {
	if key == defines.EVTUpdateIPV6 {
		if poster, ok := publisher.(*IPV6UpdatePoster); ok {
			return I.provider.Update(defines.IPTypeV6, poster.IP())
		} else {
			return fmt.Errorf("invalid ip poster")
		}
	}
	return nil
}

type configInit struct{}

func (c configInit) Events() map[eventbus.Key]eventbus.Priority {
	return map[eventbus.Key]eventbus.Priority{
		events.EventKeyInitialize: events.PriorityINI_ParseConfig + 10,
	}
}

func (c configInit) OnEvent(key eventbus.Key, _ eventbus.Publisher) error {
	if key == events.EventKeyInitialize {
		return applyConfigs()
	}
	return nil
}

func init() {
	if err := events.Global().Register(&configInit{}); err != nil {
		log.Fatal(err)
	}
}
