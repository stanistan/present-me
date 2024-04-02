package review

import (
	"bytes"
	"context"
	"html/template"

	"github.com/stanistan/veun"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var md = goldmark.New(goldmark.WithExtensions(extension.GFM))

type markdown string

func (m markdown) AsHTML(_ context.Context) (template.HTML, error) {
	var out bytes.Buffer
	if err := md.Convert([]byte(m), &out); err != nil {
		var empty template.HTML
		return empty, err
	}

	return template.HTML(out.String()), nil
}

func (m markdown) View(_ context.Context) (*veun.View, error) {
	return veun.V(m), nil
}
