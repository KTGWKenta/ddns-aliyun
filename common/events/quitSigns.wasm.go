// +build wasm

package events

import (
	"os"
	"syscall"
)

var quitSignals = []os.Signal{
	syscall.SIGINT,
	syscall.SIGQUIT,
}
