// +build !wasm

package events

import (
	"os"
	"syscall"
)

var quitSignals = []os.Signal{
	// syscall.SIGHUP,
	syscall.SIGTERM,
	syscall.SIGINT,
	syscall.SIGQUIT,
}
