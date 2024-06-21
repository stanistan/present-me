package home

import (
	"context"
	_ "embed"

	"github.com/stanistan/veun"
	t "github.com/stanistan/veun/template"
)

var (
	//go:embed home.tpl
	tpl     string
	homeTpl = t.MustParse("home", tpl)
)

type Home struct {
	Owner, Repo, Pull string
	SearchResults     veun.AsView
}

func (h Home) View(ctx context.Context) (*veun.View, error) {
	return veun.V(t.T{Tpl: homeTpl, Slots: t.Slots{"results": h.SearchResults}, Data: h}), nil
}
