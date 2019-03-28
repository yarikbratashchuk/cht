package mock

import (
	"github.com/yarikbratashchuk/cht/internal/cht"
)

// PubSub mocks pubsub client
type PubSub struct {
	pub chan cht.Message

	done chan struct{}
}

func NewPubSub() cht.PubSub {
	return &PubSub{
		pub:  make(chan cht.Message),
		done: make(chan struct{}),
	}
}

func (p *PubSub) Subscribe(sub chan<- cht.Message) {
	for {
		select {
		case <-p.done:
			return
		case m := <-p.pub:
			sub <- m
		}
	}
}

func (p *PubSub) Publish(m cht.Message) error {
	p.pub <- m
	return nil
}

func (p *PubSub) Stop() {
	close(p.done)
}
