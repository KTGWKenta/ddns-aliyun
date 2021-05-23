package events

import (
	"os"
	"os/signal"

	"github.com/kentalee/eventbus"
	"github.com/kentalee/log"
)

const (
	EventKeyPreCollect = "__preCollect__"
	EventKeyInitialize = "__initialize__"
	EventKeyAfterCheck = "__afterCheck__"

	EventKeyShutdown = "__shutdown__"
)

var global eventbus.Bus

func init() {
	global = eventbus.NewSyncBus()
}

func Global() eventbus.Bus {
	return global
}

type EventPreCollect struct{}

func (EventPreCollect) Event() eventbus.Key { return EventKeyPreCollect }

type EventInitialize struct{}

func (EventInitialize) Event() eventbus.Key { return EventKeyInitialize }

type EventAfterCheck struct{}

func (EventAfterCheck) Event() eventbus.Key { return EventKeyAfterCheck }

type EventShutdown struct{}

func (EventShutdown) Event() eventbus.Key { return EventKeyShutdown }

func Emit() (err error) {
	err = Global().Post(EventPreCollect{})
	if err != nil {
		return err
	}
	err = Global().Post(EventInitialize{})
	if err != nil {
		return err
	}
	err = Global().Post(EventAfterCheck{})
	if err != nil {
		return err
	}

	var quit = make(chan os.Signal)
	go func() {
		sig := <-quit
		log.Printf("received signal `%s`, going to shutdown!", sig.String())
		_ = Global().Post(EventShutdown{})
		os.Exit(0)
	}()
	signal.Notify(quit, quitSignals...)
	return nil
}
