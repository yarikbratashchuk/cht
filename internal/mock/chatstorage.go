package mock

import "github.com/yarikbratashchuk/cht/internal/cht"

type storage struct {
	users map[string]uint64
	rooms map[string]uint64

	messages map[uint64][]cht.Message
}

func NewChatStorage() cht.Storage {
	return storage{
		users: make(map[string]uint64, 10),
		rooms: make(map[string]uint64, 10),

		messages: make(map[uint64][]cht.Message, 10),
	}
}

func (s storage) RoomID(name string) (uint64, error) {
	id, ok := s.rooms[name]
	if !ok {
		id = newID(name)
		s.rooms[name] = id
	}
	return id, nil
}

func (s storage) UserID(name string) (uint64, error) {
	id, ok := s.users[name]
	if !ok {
		id = newID(name)
		s.users[name] = id
	}
	return id, nil
}

func (s storage) NewMessage(m cht.Message) (uint64, error) {
	mid := newID(m.Text + m.Author)
	s.messages[m.RoomID] = append(s.messages[m.RoomID], m)
	return mid, nil
}

func (s storage) LatestMsg(roomID uint64) ([]cht.Message, error) {
	msgs, ok := s.messages[roomID]
	if !ok || len(msgs) == 0 {
		return []cht.Message{}, nil
	}
	return msgs, nil
}

func newID(s string) uint64 {
	id := uint64(1)
	for _, c := range s {
		id += uint64(c)
	}
	return id
}
