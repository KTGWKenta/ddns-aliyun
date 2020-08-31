package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/MGEs/Base/workflow"
	"io/ioutil"
	"net/http"
	"time"
)

type ipResponse struct {
	Ip string `json:"ip"`
}

var newIPV4Msg = workflow.NewMask("newIPV4", "new ipv4, {{addr}}")
var newIPV6Msg = workflow.NewMask("newIPV6", "new ipv6, {{addr}}")

type worker struct {
	workflow.TaskMethods
	ticker     *time.Ticker
	stop       chan bool
	httpClient *http.Client
	ipv4       chan string
	ipv6       chan string
}

func (w *worker) SendRequest(url string, addrChan chan string) {
	defer func() {
		if err := recover(); err != nil {
			workflow.Throw(err, workflow.TlError)
		}
	}()

	var body []byte
	var result ipResponse
	var resp *http.Response
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	if req, err := http.NewRequest(
		"GET", url, nil,
	); err != nil {
		panic(errors.Wrap(err, "failed to build request"))
	} else if resp, err = w.httpClient.Do(req); err != nil {
		panic(errors.Wrap(err, "failed to execute request"))
	} else if body, err = ioutil.ReadAll(resp.Body); err != nil {
		panic(errors.Wrap(err, "failed to read body"))
	} else if err = json.Unmarshal(body, &result); err != nil {
		panic(errors.Wrap(err, "failed to parse body"))
	}
	addrChan <- result.Ip
}

func (w *worker) Looper() {
	var lastIpv4 string
	var lastIpv6 string
	for {
		select {
		case <-w.ticker.C:
			go w.SendRequest("http://ipv4.lookup.test-ipv6.com/ip/", w.ipv4)
			go w.SendRequest("http://ipv6.lookup.test-ipv6.com/ip/", w.ipv6)
		case <-w.stop:
			w.ticker.Stop()
			w.SetStatus(workflow.Task_Status_Stopped)
			return
		case ipv4 := <-w.ipv4:
			if ipv4 == lastIpv4 {
				break
			}
			lastIpv4 = ipv4
			//ipv4Addr <- ipv4
			workflow.Throw(workflow.NewSimpleThrowable(newIPV4Msg, map[string]string{
				"addr": ipv4,
			}), workflow.TlNotify)
		case ipv6 := <-w.ipv6:
			if ipv6 == lastIpv6 {
				break
			}
			lastIpv6 = ipv6
			//ipv6Addr <- ipv6
			workflow.Throw(workflow.NewSimpleThrowable(newIPV6Msg, map[string]string{
				"addr": ipv6,
			}), workflow.TlNotify)
		}
	}
}

func (w *worker) Start() error {
	w.SetStatus(workflow.Task_Status_Staring)

	w.httpClient = &http.Client{}
	w.ticker = time.NewTicker(time.Second * 3)
	w.stop = make(chan bool)
	w.ipv4 = make(chan string)
	w.ipv6 = make(chan string)

	go w.Looper()

	w.SetStatus(workflow.Task_Status_Running)
	return nil
}

func (w *worker) Stop() error {
	w.SetStatus(workflow.Task_Status_Stopping)
	w.stop <- true
	return nil
}

func initWorker() {
	if err := workflow.GlobalEvents().Subscribe(
		workflow.EVT_Lifecycle_BootUp, 0, "ipLooper",
		func(event workflow.Event) error { return event.CommitPayload(&worker{}) },
	); err != nil {
		workflow.Throw(err, workflow.TlPanic)
	}
}
