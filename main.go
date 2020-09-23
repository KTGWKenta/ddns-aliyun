package main

import "gitlab.com/MGEs/Base/workflow"

func init() {
	initDDNS()
	initWorker()
}

func main() {
	if err := workflow.Lifecycle_Start(); nil != err {
		panic(err)
	}
}
