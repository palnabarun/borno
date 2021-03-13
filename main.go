package main

import (
	"embed"
	"html/template"
	"io"
	"os"

	"github.com/palnabarun/borno/internal"
	flag "github.com/spf13/pflag"
)

//go:embed templates/index.tmpl
var fs embed.FS

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

	tmpl, err := template.ParseFS(fs, "templates/index.tmpl")
	if err != nil {
		return err
	}

	server, err := internal.NewServer(&internal.ServerOpts{Config: config, Logger: logger, Template: tmpl})

	return server.Run()
}
