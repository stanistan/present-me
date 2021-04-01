package main

import (
	"context"
	"log"
	"os"

	pm "github.com/stanistan/present-me"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("missing argument for url")
	}

	params, err := pm.ReviewParamsFromURL(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	model, err := params.Model(pm.Context{
		Ctx:    ctx,
		Client: pm.GithubClient(ctx),
	}, false)
	if err != nil {
		log.Fatal(err)
	}

	if err := model.AsMarkdown(os.Stdout, pm.AsMarkdownOptions{}); err != nil {
		log.Fatal(err)
	}
}
