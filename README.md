#checkurl


###Installation：

```
    go get -u github.com/syhlion/go-check-url
```

###Usage：
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

    //需要new 新的 resource 參數分別為 url,抓取間格時間
	resource := checkurl.NewResource("http://tw.yahoo.com", 1*time.Second, 1*time.Second)
	resource2 := checkurl.NewResource("http://www.google.com.tw", 1*time.Second, 1*time.Second)
	resources := []*checkurl.Resource{resource, resource2}

    //最後需要new 一個bot 傳入 []*Resource *log.Logger, ReadHead interface{}
	bot := checkurl.NewBot(resources, logger, moniter)
	bot.Start()

}
```
