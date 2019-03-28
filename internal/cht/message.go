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
	ID   uint64 `json:"-"`
	Type MessageType

	Text string

	Meta *ClientMeta
}

func (m *Message) String() string {
	a := color.New(color.FgCyan, color.Bold).Sprintf(m.Meta.UserNickname)

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
