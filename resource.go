package checkurl

import (
	"net/http"
	"time"
)

type Resource struct {
	url          string //設定url ex http://tw.yahoo.com
	PollInterval time.Duration
	errCount     int
	ErrTimeOut   time.Duration
}

func NewResource(u string, pollInterval time.Duration, ErrTimeout time.Duration) *Resource {
	return &Resource{u, pollInterval, 0, ErrTimeout}
}

func (r *Resource) Poll() (*http.Response, error) {
	res, err := http.Head(r.url)
	if err != nil {
		r.errCount++
		return nil, err
	}
	r.errCount = 0
	return res, nil
}

func (r *Resource) Sleep(done chan<- *Resource) {
	time.Sleep(r.PollInterval)
	done <- r
}
