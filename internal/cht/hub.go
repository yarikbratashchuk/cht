package cht

import "github.com/yarikbratashchuk/cht/fname"

type Hub struct {
	// client connections grouped by room id
	clients map[uint64]map[*Client]struct{}

	register   chan *Client
	unregister chan *Client

	publish   chan Message
	subscribe chan Message

	pubsub  *PubSub
	storage Storage

	done chan struct{}
}

func NewHub(pubsub *PubSub, storage Storage) *Hub {
	return &Hub{
		clients: make(map[uint64]map[*Client]struct{}),

		register:   make(chan *Client),
		unregister: make(chan *Client),

		publish:   make(chan Message, 100),
		subscribe: make(chan Message, 100),

		pubsub:  pubsub,
		storage: storage,

		done: make(chan struct{}),
	}
}

func (h *Hub) Run() {
	fn := fname.Current()

	log.Infof("starting %s", fn)
	defer log.Infof("shutting down %s", fn)

	go h.pubsub.Run(h.subscribe)
	defer h.pubsub.Stop()

	for {
		select {
		case client := <-h.register:
			log.Debugf("%s: register: %v", fn, client)
			_, ok := h.clients[client.rID]
			if !ok {
				h.clients[client.rID] = make(map[*Client]struct{})
			}
			h.clients[client.rID][client] = struct{}{}

			ms, err := h.storage.LatestMsg(client.rID)
			if err != nil {
				log.Errorf("%s: storage.LatestMsg: %v", fn, err)
				continue
			}
			for _, m := range ms {
				client.send <- m
			}

			h.publish <- Message{
				Type: JoinRoom,

				RoomID: client.rID,
				Author: client.nickname,
			}

		case client := <-h.unregister:
			log.Debugf("%s: unregister: %v", fn, client)
			if rClients, ok := h.clients[client.rID]; !ok {
				log.Errorf("%s: unregister: invalid client", fn)
			} else {
				delete(rClients, client)
				close(client.send)
			}

			h.publish <- Message{
				Type: LeaveRoom,

				RoomID: client.rID,
				Author: client.nickname,
			}

		case m := <-h.publish:
			log.Debugf("%s: publish: %v", fn, m)
			var err error
			if m.ID, err = h.storage.NewMessage(m); err != nil {
				log.Errorf("%s: storage.SaveMessage: %v", fn, err)
				continue
			}
			if err := h.pubsub.Publish(m); err != nil {
				log.Errorf("%s: pubsub.Publish: %v", fn, err)
				continue
			}

		case m := <-h.subscribe:
			log.Debugf("%s: subscribe: %v", fn, m)
			if m.RoomID == 0 {
				log.Errorf("%s: subscribe: m.RoomID=0", fn)
				continue
			}
			for client := range h.clients[m.RoomID] {
				client.send <- m
			}

		case <-h.done:
			return
		}
	}
}

func (h *Hub) Stop() { close(h.done) }
