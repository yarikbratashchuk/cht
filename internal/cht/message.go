package cht

import (
	"fmt"

	"github.com/fatih/color"
)

type MessageType int

const (
	JoinRoom MessageType = 1 + iota
	LeaveRoom
	Text
)

type Message struct {
	Type MessageType

	ID     uint64 `json:"-"`
	UserID uint64 `json:"-"`
	RoomID uint64 `json:"-"`

	Author string `json:",omitempty"`
	Text   string
}

func (m *Message) String() string {
	a := color.New(color.FgCyan, color.Bold).Sprintf(m.Author)

	switch m.Type {
	case JoinRoom:
		return fmt.Sprintf("%s connected", a)
	case LeaveRoom:
		return fmt.Sprintf("%s disconnected", a)
	case Text:
		return fmt.Sprintf("%s: %s", a, m.Text)
	}

	return fmt.Sprintf("unknown message type: %v", m.Type)
}
