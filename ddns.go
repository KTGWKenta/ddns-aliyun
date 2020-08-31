package main

import "gitlab.com/MGEs/Base/workflow"



func initDDNS() {
	if err := workflow.GlobalEvents().Subscribe(
		workflow.EVT_Lifecycle_Initialize, workflow.EventPriority_Latest, "ddns",
		func(event workflow.Event) error { return nil },
	); err != nil {
		workflow.Throw(err, workflow.TlPanic)
	}
}
