package main

import (
	"context"
	"log"
	"os"

	pm "github.com/stanistan/present-me"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("missing argument for url & config")
	}

	params, err := pm.ReviewParamsFromURL(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	config := pm.MustConfig(os.Args[2])
	g, err := pm.NewGH(config.Github)
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
