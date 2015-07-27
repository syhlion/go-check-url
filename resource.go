package main

import (
	"net/http"
	"time"
)

type Resource struct {
	url          string //設定url ex http://tw.yahoo.com
	pollInterval time.Duration
	errCount     int
	errTimeOut   time.Duration
}

type Resourcer interface {
	ReadHead(rs *http.Response, err error)
}

func (r *Resource) poll(callback Resourcer) {
	res, err := http.Head(r.url)
	callback.ReadHead(res, err)
}

func (r *Resource) sleep(done chan<- *Resource) {
	time.Sleep(r.pollInterval + r.errTimeOut*time.Duration(r.errCount))
	done <- r
}
