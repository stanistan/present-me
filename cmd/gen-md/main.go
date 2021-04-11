package main

import (
	"context"
	"os"

	"github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"

	pm "github.com/stanistan/present-me"
)

func main() {
	var c struct {
		Config pm.Config `embed:"" required:""`
		URL    string    `arg:"" required:""`
	}
	_ = kong.Parse(&c)

	params, err := pm.ReviewParamsFromURL(c.URL)
	if err != nil {
		log.Fatal(err)
	}

	g, err := pm.NewGH(c.Config.Github)
	if err != nil {
		log.Fatal(err)
	}

	model, err := params.Model(context.Background(), g, false)
	if err != nil {
		log.Fatal(err)
	}

	if err := model.AsMarkdown(os.Stdout, pm.AsMarkdownOptions{}); err != nil {
		log.Fatal(err)
	}
}
