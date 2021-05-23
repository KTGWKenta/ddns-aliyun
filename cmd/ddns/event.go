package main

import (
	"github.com/kentalee/eventbus"
	"github.com/kentalee/log"

	"github.com/kentalee/ddns/common/events"
)

// log events

type logEvents struct {
	flusher func()
}

func (s logEvents) Events() map[eventbus.Key]eventbus.Priority {
	return map[eventbus.Key]eventbus.Priority{
		events.EventKeyPreCollect: events.PriorityPRE_InitLogger,
		events.EventKeyShutdown:   events.PrioritySHU_FlushLogger,
	}
}

func (s *logEvents) OnEvent(key eventbus.Key, _ eventbus.Publisher) error {
	switch key {
	case events.EventKeyPreCollect:
		s.flusher = log.Setup("zap", log.ModeDevelop)
	case events.EventKeyShutdown:
		if s.flusher != nil {
			s.flusher()
		}
	}
	return nil
}

// shutdown events

type shutdownEvents struct{}

func (s shutdownEvents) Events() map[eventbus.Key]eventbus.Priority {
	return map[eventbus.Key]eventbus.Priority{
		events.EventKeyShutdown: events.PrioritySHU_AppShutdown,
	}
}

func (s shutdownEvents) OnEvent(key eventbus.Key, _ eventbus.Publisher) error {
	switch key {
	case events.EventKeyShutdown:
		shutdown()
	}
	return nil
}

func init() {
	if err := events.Global().Register(&logEvents{}); err != nil {
		panic(err)
	}
	if err := events.Global().Register(&shutdownEvents{}); err != nil {
		panic(err)
	}
}
