package main

import (
	"io"
	l "log"
	"net/http"
	"os"
	"os/signal"

	"github.com/btcsuite/btclog"
	flags "github.com/jessevdk/go-flags"
	"github.com/yarikbratashchuk/cht/internal/cht"
	"github.com/yarikbratashchuk/cht/internal/mock"
)

var log btclog.Logger

func fatalf(format string, params ...interface{}) {
	log.Criticalf(format, params...)
	os.Exit(1)
}

func setupLog(dest io.Writer, loglevel string) {
	logBackend := btclog.NewBackend(dest)
	lvl, _ := btclog.LevelFromString(loglevel)

	chtLog := logBackend.Logger("CHAT")
	log = logBackend.Logger("SRVR")

	chtLog.SetLevel(lvl)
	log.SetLevel(lvl)

	cht.UseLogger(chtLog)
}

func main() {
	conf, err := loadConfig()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok &&
			flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		l.Fatalf("loading config: %v\n", err)
	}

	setupLog(os.Stderr, conf.LogLevel)

	storage := mock.NewChatStorage()
	pubsub := mock.NewPubSub()
	hub := cht.NewHub(pubsub, storage)

	go hub.Run()
	defer hub.Stop()

	h := cht.NewServer(hub)

	server := &http.Server{
		Addr:    ":" + conf.Port,
		Handler: h,
	}
	go server.ListenAndServe()

	// Shutdown on SIGINT (CTRL-C).
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Infof("shutting down server...")
}
