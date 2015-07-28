package checkurl

import (
	"log"
	"net/http"
)

type Checkbot struct {
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
	ReadHead() chan<- State
}

func NewBot(r []*Resource, l *log.Logger, cb StateReader) *Checkbot {
	return &Checkbot{r, l, cb, 3}

}

func poller(in <-chan *Resource, out chan<- *Resource, state chan<- State) {
	for i := range in {
		res, err := i.Poll()
		out <- i
		state <- State{Resp: res, Err: err, Url: i.url}
	}
}

func (b *Checkbot) Start() {

	b.logs.Println("Bot Start...")
	pending := make(chan *Resource)
	complete := make(chan *Resource)

	state := b.ReadHead()

	for i := 0; i < b.Threads; i++ {
		go poller(pending, complete, state)
	}

	for _, r := range b.resources {
		pending <- r
	}

	for c := range complete {
		c.Sleep(pending)
	}

}
