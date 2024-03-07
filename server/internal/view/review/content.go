package review

import (
	"fmt"

	"github.com/stanistan/present-me/internal/api"
	"github.com/stanistan/present-me/internal/github"
	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"
)

func topBar(main, right veun.AsView) el.Div {
	return el.Div{
		el.Class("bg-gradient-to-b from-gray-800 to-black", "text-white", "font-mono", "text-sm", "py-2", "shadow"),
		el.Div{
			el.Class("flex flex-row"),
			el.Div{
				el.Class("flex-none pl-3"),
				el.A{el.Href("/"), el.Class("font-bold hover:text-pink-300"), el.Text(" / ")},
			},
			el.Div{el.Class("flex-grow")},
			el.Div{el.Class("flex-none"), el.Content{main}},
			el.Div{el.Class("flex-grow")},
			el.Div{el.Class("flex-none pr-3"), el.Content{right}},
		},
	}
}

func PageContent(p github.ReviewParamsMap, model api.Review) el.Div {
	return el.Div{
		topBar(
			// TODO: topbar needs the buttons to do slides and cards, etc
			el.Text(fmt.Sprintf("%s/%s#%s", p.Owner, p.Repo, p.Pull)),
			nil,
		),
		el.Div{
			el.Class("relative"),
			el.Div{
				el.Class("h-full"),
				el.Div{
					el.Class("gap-3"),
					el.Div{
						el.Class("pt-4"),
						el.Text("ReviewMetadataList for model"),
					},
					Card{
						Title: el.Div{
							el.Class("text-xsl font-extrabold"),
							el.Text(model.Title.Text),
						},
						Body: el.Div{
							el.Class("p-4"),
							el.Content{
								Markdown(model.Body),
							},
						},
					}.View(),
					el.MapFragment(model.Comments, func(c api.Comment, idx int) el.Div {
						return Card{
							Badge: idx + 1,
							Title: el.A{
								el.Code{
									el.Text(c.Title.Text),
								},
								el.Class("underline", "hover:no-underline"),
								el.Href(c.Title.HRef),
								el.Attr{"target", "_blank"},
							},
							Body: el.Div{
								el.Class("flex flex-col md:flex-row max-h-[95vh] bg-gray-50"),
								el.Div{
									el.Class("p-3 flex-none md:w-2/5 text-md markdown"),
									el.Content{Markdown(c.Description)},
								},
								el.Div{
									el.Class("flex-grow overflow-scroll text-sm md:border-l border-t md:border-t-0"),
									Diff(c.CodeBlock),
								},
							},
						}.View()
					}),
				},
			},
		},
	}
}

// TODO: markdown
func Markdown(in string) veun.AsView {
	return el.Div{el.Text(in)}
}

func Diff(code api.CodeBlock) el.Pre {
	return el.Pre{
		el.Class("bg-gray-100"),
		el.Code{
			el.Class(code.Language),
			el.Text(code.Content),
		},
	}
}

type Card struct {
	Badge       int
	Title, Body veun.AsView
}

func (c Card) View() el.Div {
	return el.Div{
		el.Class("m-4", "border", "border-slate-300", "rounded-xl", "overflow-hidden", "shadow"),
		el.Div{
			el.Class("w-full"),
			el.Div{
				el.Class("w-full", "text-sm", "p-3", "bg-slate-100", "border-b", "border-slate-300", "gap-1", "rounded-t-xl"),
				badge(c.Badge),
				el.Content{c.Title},
			},
			el.Content{c.Body},
		},
	}
}

func badge(i int) el.Fragment {
	if i == 0 {
		return nil
	}

	return el.Fragment{
		el.Span{
			el.Class("text-center ring bg-gray-700 text-white rounded-3xl p-1 px-2 text-xs mr-2 ring-gray-100 font-mono"),
			el.Text(fmt.Sprintf("%d", i)),
		},
	}
}
