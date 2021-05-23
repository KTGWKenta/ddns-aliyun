package internal

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/kentalee/errors"
	"github.com/kentalee/log"

	"github.com/kentalee/ddns/common/events"
)

type Worker struct {
	ticker     *time.Ticker
	stop       chan bool
	httpClient *http.Client
	ipv4       chan string
	ipv6       chan string
}

func (w *Worker) fetchIP(url string, addrChan chan string) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
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
		log.Error(errors.Note(err, zap.String("reason", "failed to build request")))
	}
	if resp, err = w.httpClient.Do(req); err != nil {
		log.Error(errors.Note(err, zap.String("reason", "failed to execute request")))
	}
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Error(errors.Note(err, zap.String("reason", "failed to read body")))
	}
	if err = json.Unmarshal(body, &result); err != nil {
		log.Error(errors.Note(err, zap.String("reason", "failed to parse body")))
	}
	addrChan <- result.Ip
}

func (w *Worker) looper() {
	var lastIpv4 string
	var lastIpv6 string
	if Config.Lookups.V4Enable {
		go w.fetchIP(Config.Lookups.V4Addr, w.ipv4)
	}
	if Config.Lookups.V6Enable {
		go w.fetchIP(Config.Lookups.V6Addr, w.ipv6)
	}
	for {
		select {
		case <-w.ticker.C:
			if Config.Lookups.V4Enable {
				go w.fetchIP(Config.Lookups.V4Addr, w.ipv4)
			}
			if Config.Lookups.V6Enable {
				go w.fetchIP(Config.Lookups.V6Addr, w.ipv6)
			}
		case <-w.stop:
			w.ticker.Stop()
			return
		case ipv4 := <-w.ipv4:
			if ipv4 == lastIpv4 {
				break
			}
			lastIpv4 = ipv4
			if err := events.Global().Post(&IPV4UpdatePoster{ip: ipv4}); err != nil {
				log.Error(err)
			}
		case ipv6 := <-w.ipv6:
			if ipv6 == lastIpv6 {
				break
			}
			lastIpv6 = ipv6
			if err := events.Global().Post(&IPV6UpdatePoster{ip: ipv6}); err != nil {
				log.Error(err)
			}
		}
	}
}

func (w *Worker) Start() error {
	w.httpClient = &http.Client{
		Timeout: Config.Lookups.Timeout,
	}
	w.ticker = time.NewTicker(Config.Lookups.Interval)
	w.stop = make(chan bool)
	w.ipv4 = make(chan string)
	w.ipv6 = make(chan string)

	w.looper()
	return nil
}

func (w *Worker) Stop() error {
	w.stop <- true
	return nil
}
