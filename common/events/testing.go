package events

import (
	"testing"

	"github.com/kentalee/log"
)

func MainInjector(m *testing.M, beforeRun, afterRun func()) {
	if err := Emit(); err != nil {
		log.Fatal(err)
	}
	if beforeRun != nil {
		beforeRun()
	}
	m.Run()
	if afterRun != nil {
		afterRun()
	}
	if err := Global().Post(EventShutdown{}); err != nil {
		log.Fatal(err)
	}
}
