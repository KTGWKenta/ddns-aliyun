package main

import "gitlab.com/MGEs/Base/workflow"

var ipv4Addr = make(chan string)
var ipv6Addr = make(chan string)

func main() {
	if err := workflow.Lifecycle_Start(); nil != err {
		panic(err)
	}
}
