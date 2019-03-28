package cht

type PubSub interface {
	Publish(Message) error
	Subscribe(chan<- Message)

	Stop()
}
