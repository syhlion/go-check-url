# httpbot

Help me use http crawler can custom everyone response & request

It base on [Gorilla Websocket](https://github.com/gorilla/websocket)

### Installation：

```
    go get -u github.com/syhlion/go-httpbot
```

### Usage：
```
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/syhlion/go-httpbot"
)

//客製化得到的Response 後續處理
type CustomResponseReader struct {
	*log.Logger
}

func (c *CustomResponseReader) Read(resp *http.Response) (res *http.Response, err error) {

    //仿造得到一個response後  在重新用得到的結果組出新的request 回傳出去
	client := &http.Client{}
	c.Println(resp.Request.URL.Host)
	if req, err := http.NewRequest("GET", "http://www.google.com", nil); err == nil {
		res, err = client.Do(req)
		c.Println(res.Request.URL.Host, err)
	}
	return

}

//最後接收回來的結果
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

var logger *log.Logger

func main() {

	logger = log.New(os.Stdout, "logger:", log.Ldate|log.Ltime)
	moniter := &Moniter{logger}
	req, err := http.NewRequest("GET", "http://www.googlw.com", nil)
	if err != nil {
		panic(err)
	}
    //客制response 鍊 可以不斷串接
	quene := []CustomResponseReader{&CustomResponseReader{logger}}

	resource := httpbot.NewResource(req, 1*time.Second,quene, nil)
	resources := []*httpbot.Resource{resource}

	//最後需要new 一個bot 傳入 []*Resource *log.Logger, StateReader interface{}
	bot := httpbot.NewBot(resources, logger, moniter)
	bot.Start()

}
```
