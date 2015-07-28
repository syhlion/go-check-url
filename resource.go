package checkurl

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

func NewResource(u string) *Resource {
	return &Resource{u, 60 * time.Second, 0, 10 * time.Second}
}

func (r *Resource) Poll() (*http.Response, error) {
	res, err := http.Head(r.url)
	if err != nil {
		r.errCount++
	} else {
		r.errCount = 0
	}
	return res, err
}

func (r *Resource) Sleep(done chan<- *Resource) {
	time.Sleep(r.pollInterval + r.errTimeOut*time.Duration(r.errCount))
	done <- r
}
