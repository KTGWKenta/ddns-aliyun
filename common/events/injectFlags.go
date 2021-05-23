package events

import (
	"flag"

	"github.com/kentalee/eventbus"
)

type flagEvent struct{}

func (f flagEvent) Events() map[eventbus.Key]eventbus.Priority {
	return map[eventbus.Key]eventbus.Priority{
		EventKeyInitialize: PriorityINI_ParseFlag,
	}
}

func (f flagEvent) OnEvent(key eventbus.Key, _ eventbus.Publisher) error {
	switch key {
	case EventKeyInitialize:
		flag.Parse()
	}
	return nil
}

func init() {
	if err := Global().Register(&flagEvent{}); err != nil {
		panic(err)
	}
}
