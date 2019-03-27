// The cht is chat cli
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
	c "github.com/yarikbratashchuk/cht/internal/cli"
)

const (
	defaultRoom     = "random"
	defaultNickname = "noname"
	defaultServer   = "127.0.0.1:9090"
)

func main() {
	app := cli.NewApp()
	app.Name = "cht"
	app.Version = "1.0.0"
	app.Usage = "chat cli"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  c.RoomFlag,
			Value: defaultRoom,
			Usage: "room to join",
		},
		cli.StringFlag{
			Name:  c.ServerFlag,
			Value: defaultServer,
			Usage: "server to connect",
		},
		cli.StringFlag{
			Name:  c.NicknameFlag,
			Value: defaultNickname,
			Usage: "your nickname",
		},
	}
	app.Commands = []cli.Command{
		c.StartCmd,
	}

	fmt.Print("\n")
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
	fmt.Print("\n")
}
