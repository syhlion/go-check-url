package httpbot

import (
	"net/http"
	"time"
)

type Resource struct {
	HttpRequest *http.Request

	// Set Poll Interval
	PollInterval time.Duration

	respQuene  []ResponseReader
	httpClient *http.Client
}

//ResponseReader
//可以每個ResponseReader 都可以接受到後客製化新的request重新送出
type ResponseReader interface {
	Read(resp *http.Response) (*http.Response, error)
}

func NewResource(httpreq *http.Request, pollInterval time.Duration, respQuene []ResponseReader, client *http.Client) *Resource {
	if client == nil {
		client = &http.Client{}
	}
	return &Resource{httpreq, pollInterval, respQuene, client}
}

func (r *Resource) Get() (res *http.Response, err error) {
	res, err = r.httpClient.Do(r.HttpRequest)
	for _, cb := range r.respQuene {
		res, err = cb.Read(res)
	}

	return
}

func (r *Resource) sleep(done chan<- *Resource) {
	time.Sleep(r.PollInterval)
	done <- r
}
