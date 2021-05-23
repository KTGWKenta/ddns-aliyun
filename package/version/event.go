package version

import (
	"flag"
	"os"

	"github.com/kentalee/eventbus"

	"github.com/kentalee/ddns/common/events"
)

type flagSubscriber struct {
	doCheck bool
}

func (f flagSubscriber) Events() map[eventbus.Key]eventbus.Priority {
	return map[eventbus.Key]eventbus.Priority{
		events.EventKeyInitialize: events.PriorityINI_CheckFlag,
	}
}

func (f *flagSubscriber) OnEvent(key eventbus.Key, _ eventbus.Publisher) error {
	switch key {
	case events.EventKeyInitialize:
		if f.doCheck {
			Print()
			os.Exit(0)
		}
	}
	return nil
}

func init() {
	subscriber := flagSubscriber{}
	flag.BoolVar(&subscriber.doCheck, "version", false, "check the version")
	if err := eventbus.GetDefault().Register(&subscriber); err != nil {
		panic(err)
	}
}
