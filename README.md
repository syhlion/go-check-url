# checkurl

練習go concurrent，此套件可以拿同時查詢關注的網頁是否存活

### Installation：

```
    go get -u github.com/syhlion/go-check-url
```

### Usage：
```
package main

import (
	"flag"
	"github.com/syhlion/go-check-url"
	"log"
	"os"
	"time"
)

type Moniter struct {
	*log.Logger
}

//實作 ReadHead interface 當作 callback
func (m *Moniter) ReadHead() chan<- checkurl.State {
	updates := make(chan checkurl.State)
	go func() {
		for {
			select {
			case s := <-updates:
				if s.Err != nil {
					m.Println(s.Err)
				} else {
					m.Println(s.Url, s.Resp.Status)
				}
			}
		}
	}()

	return updates
}
//可用cmd傳入 要log的位置
var logDir = flag.String("logDir", "", "log dir")
var logger *log.Logger

func main() {
	flag.Parse()
	f, err := os.OpenFile(*logDir, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0775)
	if err == nil {
		logger = log.New(f, "logger:", log.Ldate|log.Ltime)
	} else {

        //如果沒有餵入參數使用系統的stdout
		logger = log.New(os.Stdout, "logger:", log.Ldate|log.Ltime)
	}
	moniter := &Moniter{logger}
    18     //  req, err := http.NewRequest("GET", u, nil)
    19     //  if err != nil {
    20     //      panic(err)
    21     //  }
    22     //  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    23     //  req.Header.Set("Accept-Encoding", "gzip,deflate")
    24     //  req.Header.Set("Connection", "keep-alive")
    25     //  req.Header.Set("Accept-Language", "en-US,en;q=0.5")
    26     //  req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:39.0) Gecko/20100101 Firefox/39.0")

    //需要new 新的 resource 參數分別為 url,抓取間格時間
	resource := checkurl.NewResource("http://tw.yahoo.com", 1*time.Second)
	resource2 := checkurl.NewResource("http://www.google.com.tw", 1*time.Second)
	resources := []*checkurl.Resource{resource, resource2}

    //最後需要new 一個bot 傳入 []*Resource *log.Logger, ReadHead interface{}
	bot := checkurl.NewBot(resources, logger, moniter)
	bot.Start()

}
```
