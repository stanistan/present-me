package review

import (
	"fmt"

	"github.com/stanistan/present-me/internal/api"
	"github.com/stanistan/present-me/internal/github"
	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"
)

func SourcesFragment(
	p github.ReviewParamsMap, m api.Review, sources []api.SourceProvider,
) el.Fragment {
	content := el.MapFragment(
		sources,
		func(s api.SourceProvider, _ int) el.Param {
			label := s.Label()
			href := s.Link()
			return el.Li{
				el.A{
					el.Class(
						"block",
						"px-3", "py-2", "border-b",
						"text-sm", "bg-white hover:bg-pink-50",
						"font-mono", "hover:text-indigo-900",
					),
					el.Href(href),
					el.Div{
						el.Class("underline"),
						el.Text(label),
					},
				},
			}
		},
	)

	if len(content) == 0 {
		content = el.Fragment{
			el.Div{
				el.Class("text-center", "p-2"),
				el.Text("none"),
			},
		}
	}

	return el.Fragment{
		Card{
			Title: el.Div{
				el.Class("text-md", "text-center"),
				el.Span{
					el.Class("font-bold"),
					el.Text(fmt.Sprintf("%s/%s#%s", p.Owner, p.Repo, p.Pull)),
				},
			},
			Body: el.Div{
				el.Ol{
					el.Class("list-decimal"),
					content,
				},
			},
		}.Render(),
	}
}

func SourcesList(p github.ReviewParamsMap, m api.Review, sources []api.SourceProvider) veun.AsView {
	return Layout(p, m, LayoutParams{
		Content: el.Fragment{
			el.Class("max-w-96", "p-3", "mx-auto"),
			SourcesFragment(p, m, sources),
		},
	})
}
