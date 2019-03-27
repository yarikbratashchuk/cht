package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"github.com/apex/log"
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

	fmt.Println("connecting...")

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

	cx, cancel := context.WithCancel(context.Background())
	go read(cx, conn, os.Stdout)
	go write(cx, conn, os.Stdin)

	d := make(chan os.Signal, 1)
	signal.Notify(d, os.Interrupt)
	<-d
	cancel()

	fmt.Println("disconnecting...")

	return nil
}

type messageWriter interface {
	WriteMessage(messageType int, data []byte) error
}

func write(ctx context.Context, w messageWriter, r io.Reader) {
	messageChan := make(chan cht.Message)
	go func() {
		for {
			messageChan <- readMessage(r)
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
				log.Errorf("marshal: %v", err)
				return
			}
			err = w.WriteMessage(
				websocket.TextMessage,
				data,
			)
			if err != nil {
				log.Errorf("write: %v", err)
				return
			}
		}
	}
}

type messageReader interface {
	ReadMessage() (messageType int, p []byte, err error)
}

func read(ctx context.Context, r messageReader, w io.Writer) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, message, err := r.ReadMessage()
			if err != nil {
				log.Errorf("recv: %v", err)
				return
			}
			var m cht.Message
			err = json.Unmarshal(message, &m)
			if err != nil {
				log.Errorf("unmarshal: %v", err)
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
		log.Errorf("write close: %v", err)
	}
}
