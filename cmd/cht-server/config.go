package main

import (
	flags "github.com/jessevdk/go-flags"
)

type config struct {
	Port string `short:"p" long:"port" description:"port to listen on"`

	LogLevel string `long:"loglevel" description:"log level for all subsystems {trace, debug, info, error, critical}"`
}

var defconf = config{
	Port: "9090",

	LogLevel: "info",
}

func loadConfig() (*config, error) {
	conf := defconf
	_, err := flags.Parse(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
