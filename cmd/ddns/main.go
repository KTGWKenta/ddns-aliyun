package main

import (
	"github.com/kentalee/log"

	_ "github.com/kentalee/ddns-aliyun/common/config"
	"github.com/kentalee/ddns-aliyun/common/events"
	"github.com/kentalee/ddns-aliyun/internal"
	_ "github.com/kentalee/ddns-aliyun/package/version"
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
