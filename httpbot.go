package httpbot

import (
	"io/ioutil"
	"log"
	"net/http"
)

type Bot struct {
	resources []*Resource
	logs      *log.Logger
	StateReader
	Threads int
}

type State struct {
	Url  string
	Resp *http.Response
	Err  error
}

//自行實作 StateReader
type StateReader interface {
	Read() chan<- State
}

func NewBot(r []*Resource, l *log.Logger, cb StateReader) *Bot {
	return &Bot{r, l, cb, 3}

}

func ReadHtml(resp *http.Response) (ret string) {
	defer resp.Body.Close()

	if body, err := ioutil.ReadAll(resp.Body); err == nil {
		ret = string(body)
	}
	return

}

func poller(in <-chan *Resource, out chan<- *Resource, state chan<- State) {
	for i := range in {
		res, err := i.poll()
		out <- i
		state <- State{Resp: res, Err: err, Url: i.HttpRequest.URL.String()}
	}
}

func (b *Bot) Start() {

	b.logs.Println("Bot Start...")
	pending := make(chan *Resource)
	complete := make(chan *Resource)

	state := b.Read()

	for i := 0; i < b.Threads; i++ {
		go poller(pending, complete, state)
	}

	for _, r := range b.resources {
		pending <- r
	}

	for c := range complete {
		c.sleep(pending)
	}

}
