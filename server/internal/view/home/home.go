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
	Owner, Repo, PR string
}

func (h Home) View(ctx context.Context) (*veun.View, error) {
	return veun.V(t.T{Tpl: homeTpl, Data: h}), nil
}
