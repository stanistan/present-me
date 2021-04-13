package main

import (
	"context"
	"os"

	"github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"

	pm "github.com/stanistan/present-me"
)

var c struct {
	Config pm.Config `embed:"" required:""`
	URL    string    `arg:"" required:""`
}

func main() {
	_ = kong.Parse(&c)

	g, err := pm.NewGH(c.Config.Github)
	if err != nil {
		log.Fatal(err)
	}

	params, err := pm.ReviewParamsFromURL(c.URL)
	if err != nil {
		log.Fatal(err)
	}

	model, err := params.Model(context.Background(), g)
	if err != nil {
		log.Fatal(err)
	}

	if err := model.AsMarkdown(os.Stdout, pm.AsMarkdownOptions{}); err != nil {
		log.Fatal(err)
	}
}
