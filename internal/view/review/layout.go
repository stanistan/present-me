package review

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/stanistan/present-me/internal/api"
	"github.com/stanistan/present-me/internal/github"
	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"
)

type TopBar struct {
	PlayButton   veun.AsView
	ViewSelector veun.AsView
}

type LayoutParams struct {
	TopBar  TopBar
	Content el.Fragment
}

func Layout(p github.ReviewParamsMap, model api.Review, layout LayoutParams) veun.AsView {
	log.Info().Any("params", p).Msg("in layout")
	return el.Div{

		topBar(
			veun.Views{
				el.Span{
					el.Class("px-3"),
					el.Text(fmt.Sprintf("%s/%s#%s", p.Owner, p.Repo, p.Pull)),
				},
				layout.TopBar.ViewSelector,
			},
			layout.TopBar.PlayButton,
		),

		el.Div{
			layout.Content,
		},
	}
}
