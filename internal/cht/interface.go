package cht

type Storage interface {
	RoomID(name string) (uint64, error)
	UserID(name string) (uint64, error)

	NewMessage(message Message) (uint64, error)
	LatestMsg(roomID uint64) ([]Message, error)
}

type PubSub interface {
	Publish(Message) error
	Subscribe(chan<- Message)

	Stop()
}
