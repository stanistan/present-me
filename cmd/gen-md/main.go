package main

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/github"
	"github.com/stanistan/present-me"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("missing argument for url")
	}

	params, err := crap.ReviewParamsFromURL(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	model, err := params.Model(crap.Context{
		Ctx:    context.Background(),
		Client: github.NewClient(nil),
	}, false)
	if err != nil {
		log.Fatal(err)
	}

	if err := model.AsMarkdown(os.Stdout, crap.AsMarkdownOptions{}); err != nil {
		log.Fatal(err)
	}
}
