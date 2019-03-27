package cht

type Message struct {
	ID     uint64 `json:"-"`
	UserID uint64 `json:"-"`
	RoomID uint64 `json:"-"`

	Author string `json:",omitempty"`
	Text   string
}

type Storage interface {
	RoomID(name string) (uint64, error)
	UserID(name string) (uint64, error)

	NewMessage(message Message) (uint64, error)
	LatestMsg(roomID uint64) ([]Message, error)
}
