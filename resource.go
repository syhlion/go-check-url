package checkurl

import (
	"net/http"
	"time"
)

type Resource struct {
	url          string        //設定url ex http://tw.yahoo.com
	PollInterval time.Duration //設定抓取間隔時間
	errCount     int
}

func NewResource(u string, pollInterval time.Duration) *Resource {
	return &Resource{u, pollInterval, 0}
}

func (r *Resource) Poll() (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", r.url, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept-Encoding", "gzip,deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:39.0) Gecko/20100101 Firefox/39.0")
	res, err := client.Do(req)
	//res, err := http.Head(r.url)
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
