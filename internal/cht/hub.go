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

			rID := client.meta.RoomID

			_, ok := h.clients[rID]
			if !ok {
				h.clients[rID] = make(map[*Client]struct{})
			}
			h.clients[rID][client] = struct{}{}

			ms, err := h.storage.LatestMsg(rID)
			if err != nil {
				log.Errorf("%s: storage.LatestMsg: %v", fn, err)
				continue
			}

			for _, m := range ms {
				client.send <- m
			}

			h.publish <- Message{
				Type: JoinRoom,

				Meta: client.meta,
			}

		case client := <-h.unregister:
			log.Debugf("%s: unregister: %v", fn, client)

			rID := client.meta.RoomID

			if roomClients, ok := h.clients[rID]; !ok {
				log.Errorf("%s: unregister: invalid client", fn)
			} else {
				delete(roomClients, client)
				close(client.send)
			}

			h.publish <- Message{
				Type: LeaveRoom,

				Meta: client.meta,
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

			if m.Meta.RoomID == 0 {
				log.Errorf("%s: subscribe: m.RoomID=0", fn)
				continue
			}

			for client := range h.clients[m.Meta.RoomID] {
				client.send <- m
			}

		case <-h.done:
			return
		}
	}
}

func (h *Hub) Stop() { close(h.done) }
