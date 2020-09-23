package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/KTGWKenta/ddns-aliyun/config"
	"github.com/KTGWKenta/ddns-aliyun/defines"
	"github.com/pkg/errors"
	"gitlab.com/MGEs/Base/workflow"
)

type ipResponse struct {
	Ip string `json:"ip"`
}

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
	var err error
	var req *http.Request
	if req, err = http.NewRequest(
		"GET", url, nil,
	); err != nil {
		workflow.Throw(errors.Wrap(err, "failed to build request"), workflow.TlPanic)
	}
	if resp, err = w.httpClient.Do(req); err != nil {
		workflow.Throw(errors.Wrap(err, "failed to execute request"), workflow.TlPanic)
	}
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		workflow.Throw(errors.Wrap(err, "failed to read body"), workflow.TlPanic)
	}
	if err = json.Unmarshal(body, &result); err != nil {
		workflow.Throw(errors.Wrap(err, "failed to parse body"), workflow.TlPanic)
	}
	addrChan <- result.Ip
}

func (w *worker) Looper() {
	var lastIpv4 string
	var lastIpv6 string
	go w.SendRequest(config.Config.Lookups.V4Addr, w.ipv4)
	go w.SendRequest(config.Config.Lookups.V6Addr, w.ipv6)
	for {
		select {
		case <-w.ticker.C:
			go w.SendRequest(config.Config.Lookups.V4Addr, w.ipv4)
			go w.SendRequest(config.Config.Lookups.V6Addr, w.ipv6)
		case <-w.stop:
			w.ticker.Stop()
			w.SetStatus(workflow.Task_Status_Stopped)
			return
		case ipv4 := <-w.ipv4:
			if ipv4 == lastIpv4 {
				break
			}
			lastIpv4 = ipv4
			if err := workflow.GlobalEvents().Publish(
				defines.EVTUpdateIPV4, workflow.NewEventAction(nil, ipv4),
			); err != nil {
				workflow.Throw(err, workflow.TlError)
			}
		case ipv6 := <-w.ipv6:
			if ipv6 == lastIpv6 {
				break
			}
			lastIpv6 = ipv6
			if err := workflow.GlobalEvents().Publish(
				defines.EVTUpdateIPV6, workflow.NewEventAction(nil, ipv6),
			); err != nil {
				workflow.Throw(err, workflow.TlError)
			}
		}
	}
}

func (w *worker) Start() error {
	w.SetStatus(workflow.Task_Status_Staring)

	w.httpClient = &http.Client{}
	w.ticker = time.NewTicker(time.Minute * 3)
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
