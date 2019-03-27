package mock

import "github.com/yarikbratashchuk/cht/internal/cht"

type storage struct{}

func NewChatStorage() cht.Storage {
	return storage{}
}

func (s storage) RoomID(name string) (uint64, error) {
	return 1, nil
}

func (s storage) UserID(name string) (uint64, error) {
	c := uint64(0)
	for _, ch := range name {
		c += uint64(ch)
	}
	return c, nil
}

func (s storage) NewMessage(message cht.Message) (uint64, error) {
	c := uint64(1)
	for _, ch := range message.Text {
		c += uint64(ch)
	}
	return c, nil
}

func (s storage) LatestMsg(roomID uint64) ([]cht.Message, error) {
	return []cht.Message{
		cht.Message{
			Author: "Jarvis",
			Text:   "papa can you hear me",
		},
		cht.Message{
			Author: "Caoimhe",
			Text:   "hop hey lalaley",
		},
	}, nil
}
