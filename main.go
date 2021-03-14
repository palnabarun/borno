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
	flag.StringVarP(&configOpts.AssetsRoot, "assets-dir", "", "assets", "path to the root of assets (DEFAULT: assets in current directory")
	flag.Parse()

	logger := internal.NewLogger(&internal.LoggerOpts{Out: out})
	configOpts.Logger = logger

	config, err := internal.ParseTalksFromConfig(configOpts)
	if err != nil {
		return err
	}

	config.Talks = internal.ProcessTalks(config.Talks)
	configOpts.Config = config

	server, err := internal.NewServer(configOpts)
	if err != nil {
		return err
	}

	return server.Run()
}
