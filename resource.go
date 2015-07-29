package httpbot

import (
	"net/http"
	"time"
)

type Resource struct {
	HttpRequest *http.Request

	// Set Poll Interval
	PollInterval time.Duration
	httpClient   *http.Client
}

func NewResource(httpreq *http.Request, pollInterval time.Duration, client *http.Client) *Resource {
	if client == nil {
		client = &http.Client{}
	}
	return &Resource{httpreq, pollInterval, client}
}

func (r *Resource) poll() (*http.Response, error) {

	res, err := r.httpClient.Do(r.HttpRequest)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Resource) sleep(done chan<- *Resource) {
	time.Sleep(r.PollInterval)
	done <- r
}
