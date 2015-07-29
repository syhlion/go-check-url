package httpbot

import (
	"net/http"
	"time"
)

type Resource struct {
	HttpRequest *http.Request

	// Set Poll Interval
	PollInterval time.Duration

	respQuene  *ResponseQuene
	httpClient *http.Client
}

type ResponseReader interface {
	Read(resp *http.Response) (*http.Response, error)
}

func NewResource(httpreq *http.Request, pollInterval time.Duration, respQuene *ResponseQuene, client *http.Client) *Resource {
	if client == nil {
		client = &http.Client{}
	}
	return &Resource{httpreq, pollInterval, respQuene, client}
}

func (r *Resource) poll() (*http.Response, error) {
	res, err := r.httpClient.Do(r.HttpRequest)
	if err != nil {
		return nil, err
	}
	for i := 0; i < r.respQuene.size; i++ {
		res, err = r.respQuene.Pop().Read(res)
	}

	return res, nil
}

func (r *Resource) sleep(done chan<- *Resource) {
	time.Sleep(r.PollInterval)
	done <- r
}
