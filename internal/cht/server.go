package cht

import (
	"fmt"
	"net/http"

	"github.com/yarikbratashchuk/cht/fname"
)

const (
	AuthHeader     = "X-Auth-Token"
	RoomHeader     = "X-Room"
	NicknameHeader = "X-Nickname"

	// TODO: remove when jwt based authentication is ready
	MockJWT = "welcomeaboard"
)

type Server struct {
	hub *Hub
}

func NewServer(hub *Hub) http.Handler {
	return Server{hub}
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fn := fname.Current()

	// TODO: JWT based authentication mechanism + TLS
	if r.Header.Get(AuthHeader) != MockJWT {
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}

	room := r.Header.Get(RoomHeader)
	nickname := r.Header.Get(NicknameHeader)

	rID, err := s.hub.storage.RoomID(room)
	if err != nil {
		log.Errorf("%s: storage.RoomID: %v", fn, err)
		write500(w)
		return
	}

	uID, err := s.hub.storage.UserID(nickname)
	if err != nil {
		log.Errorf("%s: storage.UserID: %v", fn, err)
		write500(w)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Sprintf("upgrade: %s", err), 500)
		return
	}

	metadata := &ClientMeta{
		RoomID:       rID,
		UserID:       uID,
		UserNickname: nickname,
	}

	client := &Client{
		meta: metadata,

		hub:  s.hub,
		conn: conn,

		send: make(chan Message, 10),
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func write500(w http.ResponseWriter) {
	http.Error(
		w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}
