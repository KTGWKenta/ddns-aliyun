package main

import (
	"github.com/kentalee/log"

	_ "github.com/kentalee/ddns/common/config"
	"github.com/kentalee/ddns/common/events"
	"github.com/kentalee/ddns/internal"
	_ "github.com/kentalee/ddns/package/version"
)

var w = internal.Worker{}

func main() {
	if err := events.Emit(); err != nil {
		log.Fatal(err)
	}
	if err := w.Start(); err != nil {
		log.Fatal(err)
	}
}

func shutdown() {
	if err := w.Stop(); err != nil {
		log.Fatal(err)
	}
}
