package layout

import (
	"context"
	_ "embed"

	"github.com/stanistan/veun"
	t "github.com/stanistan/veun/template"
)

type Version struct {
	URL string
	SHA string
}

type Params struct {
	Title    string
	Version  Version
	JSFiles  []string
	CSSFiles []string
}

func (p Params) maybeOverrideWith(data any) Params {
	// TODO: do the duck-typing
	return p
}

type layout struct {
	Body   veun.AsView
	Params Params
}

func Layout(data Params, body veun.AsView) veun.AsView {
	return layout{
		Body:   body,
		Params: data.maybeOverrideWith(body),
	}
}

var (
	//go:embed layout.tpl
	tpl      string
	template = t.MustParse("layout", tpl)
)

func (l layout) View(ctx context.Context) (*veun.View, error) {
	return veun.V(t.Template{
		Tpl:   template,
		Slots: t.Slots{"body": l.Body},
		Data:  l.Params,
	}), nil
}
