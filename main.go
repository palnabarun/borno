package main

import (
	"io"
	"os"

	"github.com/palnabarun/borno/internal"
	flag "github.com/spf13/pflag"
)

func main() {
	if err := run(os.Stdout); err != nil {
		panic(err)
	}
}

func run(out io.Writer) error {
	configOpts := &internal.ConfigOpts{}

	flag.StringVarP(&configOpts.ConfigFile, "config", "c", "borno.yaml", "path to the borno config (DEFAULT: borno.yaml in current directory")
	flag.Parse()

	logger := internal.NewLogger(&internal.LoggerOpts{Out: out})

	config, err := internal.ParseTalksFromConfig(configOpts)
	if err != nil {
		return err
	}

	config.Talks = internal.ProcessTalks(config.Talks)

	server, err := internal.NewServer(&internal.ServerOpts{Config: config, Logger: logger})
	if err != nil {
		return err
	}

	return server.Run()
}
