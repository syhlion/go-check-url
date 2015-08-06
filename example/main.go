package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/syhlion/go-httpbot"
)

type CustomResponseReader struct {
	*log.Logger
}

func (c *CustomResponseReader) Read(resp *http.Response) (res *http.Response, err error) {
	client := &http.Client{}
	c.Println(resp.Request.URL.Host)
	//req, _ := http.NewRequest("GET", "http://tw.yahoo.com", nil)
	if req, err := http.NewRequest("GET", "http://www.google.com", nil); err == nil {
		res, err = client.Do(req)
		c.Println(res.Request.URL.Host, err)
	}
	return

}

type Moniter struct {
	*log.Logger
}

func (m *Moniter) Read() chan<- httpbot.State {
	updates := make(chan httpbot.State)
	go func() {
		for {
			select {
			case s := <-updates:
				if s.Err != nil {
					m.Println(s.Err)
				} else {
					//html := httpbot.ReadHtml(s.Resp)
					m.Println(s.Url, s.Resp.Status)
				}
			}
		}
	}()

	return updates
}

var logDir = flag.String("logDir", "", "log dir")
var logger *log.Logger

func main() {
	flag.Parse()
	if f, err := os.OpenFile(*logDir, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0775); err == nil {
		logger = log.New(f, "logger:", log.Ldate|log.Ltime)
	} else {
		logger = log.New(os.Stdout, "logger:", log.Ldate|log.Ltime)
	}
	moniter := &Moniter{logger}
	req, err := http.NewRequest("GET", "http://tw.yahoo.com", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept-Encoding", "gzip,deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:39.0) Gecko/20100101 Firefox/39.0")

	a := CustomResponseReader{logger}

	resource2 := httpbot.NewResource(req, 1*time.Second, []httpbot.ResponseReader{&a}, nil)
	resources := []*httpbot.Resource{resource2}

	//最後需要new 一個bot 傳入 []*Resource *log.Logger, ReadHead interface{}
	bot := httpbot.NewBot(resources, logger, moniter)
	bot.Start()

}
