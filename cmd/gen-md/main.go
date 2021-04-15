package main

import (
	"context"
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"

	pm "github.com/stanistan/present-me"
)

var c struct {
	Config pm.Config `embed:"" required:""`
	URL    string    `arg:"" required:""`
}

func main() {
	_ = kong.Parse(&c)
	c.Config.Configure()
	if err := realMain(); err != nil {
		log.Fatal().Err(err)
	}
}

func realMain() error {
	g, err := pm.NewGH(c.Config.Github)
	if err != nil {
		return err
	}

	params, err := pm.ReviewParamsFromURL(c.URL)
	if err != nil {
		return err
	}

	model, err := params.Model(context.Background(), g)
	if err != nil {
		return err
	}

	if err := model.AsMarkdown(os.Stdout, pm.AsMarkdownOptions{}); err != nil {
		return err
	}

	return nil
}
