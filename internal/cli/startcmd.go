package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/urfave/cli"
	"github.com/yarikbratashchuk/cht/internal/cht"
)

var StartCmd = cli.Command{
	Name:        "start",
	Usage:       "join the chat room and start sending and receiving messages",
	ArgsUsage:   "start",
	Description: "join the chat room and start sending and receiving messages",
	Action:      start,
}

func start(ctx *cli.Context) error {
	serverHost := ctx.GlobalString(ServerFlag)
	room := ctx.GlobalString(RoomFlag)
	nickname := ctx.GlobalString(NicknameFlag)

	fmt.Printf("connecting to room %s...\n", room)

	header := make(http.Header, 3)
	// TODO: JWT based authentication
	header.Add(cht.AuthHeader, cht.MockJWT)
	header.Add(cht.RoomHeader, room)
	header.Add(cht.NicknameHeader, nickname)

	u := url.URL{Scheme: "ws", Host: serverHost}

	// TODO: TLS
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		return err
	}
	defer conn.Close()

	rCtx, cancelRead := context.WithCancel(context.Background())
	go read(rCtx, conn, os.Stdout)

	wCtx, cancelWrite := context.WithCancel(context.Background())
	go write(wCtx, conn, os.Stdin)

	d := make(chan os.Signal, 1)
	signal.Notify(d, os.Interrupt)
	<-d

	cancelRead()
	cancelWrite()

	fmt.Println("\rdisconnecting...")

	return nil
}

type messageWriter interface {
	WriteMessage(messageType int, data []byte) error
}

func write(ctx context.Context, w messageWriter, r io.Reader) {
	messageChan := make(chan cht.Message)

	ticker := time.NewTicker(50 * time.Millisecond)
	defer func() { ticker.Stop() }()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				messageChan <- readMessage(r)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			shutdownConn(w)
			return
		case message := <-messageChan:
			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("marshal: %v", err)
				return
			}
			err = w.WriteMessage(
				websocket.TextMessage,
				data,
			)
			if err != nil {
				log.Printf("write: %v", err)
				return
			}
		}
	}
}

type messageReader interface {
	ReadMessage() (messageType int, p []byte, err error)
}

func read(ctx context.Context, r messageReader, w io.Writer) {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer func() { ticker.Stop() }()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, message, err := r.ReadMessage()
			if err != nil {
				log.Printf("recv: %v", err)
				return
			}
			var m cht.Message
			err = json.Unmarshal(message, &m)
			if err != nil {
				log.Printf("unmarshal: %v", err)
				return
			}

			printMessage(w, m)
		}
	}
}

func shutdownConn(w messageWriter) {
	err := w.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(
			websocket.CloseNormalClosure,
			"",
		),
	)
	if err != nil {
		log.Printf("write close: %v", err)
	}
}
