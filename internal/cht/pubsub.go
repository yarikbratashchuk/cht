package cht

import "github.com/yarikbratashchuk/cht/fname"

// PubSub mocks pubsub client
type PubSub struct {
	pub chan Message

	done chan struct{}
}

func NewPubSub() *PubSub {
	return &PubSub{
		pub:  make(chan Message),
		done: make(chan struct{}),
	}
}

func (p *PubSub) Run(sub chan<- Message) {
	fn := fname.Current()

	for {
		select {
		case <-p.done:
			log.Debugf("%s: shutting down...", fn)
			return
		case m := <-p.pub:
			sub <- m
		}
	}
}

func (p *PubSub) Publish(m Message) error {
	p.pub <- m
	return nil
}

func (p *PubSub) Stop() {
	close(p.done)
}
