package review

import (
	"fmt"

	"github.com/stanistan/present-me/internal/api"
	"github.com/stanistan/present-me/internal/github"
	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"
)

type LayoutParams struct {
	PlayButton veun.AsView
	Content    el.Fragment
}

func Layout(p github.ReviewParamsMap, model api.Review, layout LayoutParams) veun.AsView {
	return el.Div{

		topBar(
			veun.Views{
				el.Span{
					el.Class("px-3"),
					el.Text(fmt.Sprintf("%s/%s#%s", p.Owner, p.Repo, p.Pull)),
				},
				el.Div{
					el.Class(
						"inline-block px-2 py-1",
						"bg-slate-50 shadow-inner rounded-sm",
						"text-black text-xs",
					),
					el.Fragment{
						reviewLink(p, "cards"), el.Text(" | "), reviewLink(p, "slides"),
					},
				},
			},
			layout.PlayButton,
		),

		el.Div{
			layout.Content,
		},
	}
}
