package review

import (
	_ "embed"
	"fmt"

	"github.com/stanistan/present-me/internal/api"
	"github.com/stanistan/present-me/internal/github"
	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"
)

func slide(toShow int) func(int) el.AttrFunc {
	return func(idx int) el.AttrFunc {
		cl := "hidden"
		if idx == toShow {
			cl = "visible"
		}
		return el.Class("slide", fmt.Sprintf("slide-%d", idx), cl)
	}
}

func SlideContent(p github.ReviewParamsMap, model api.Review) veun.AsView {
	slide := slide(0)

	return Layout(p, model, LayoutParams{
		TopBar: TopBar{
			PlayButton: el.Button{
				el.Class("text-xs px-2"),
				el.ID("play-full-screen"),
				el.Text("▶️"),
			},
			ViewSelector: viewSelector(p),
		},
		Content: el.Fragment{
			el.Class("relative h-[95vh]"),

			el.Div{
				el.Class("bg-white flex flex-col h-full"),
				el.ID("slideshow"),

				el.Div{el.Class("flex-grow")},
				el.Div{el.Class("flex-0 max-w-[2200px] mx-auto"),

					// slide-0
					el.Div{
						slide(0),
						el.Div{el.Class("text-6xl font-extrabold text-center"),
							GradientText(model.Title.Text),
						},
						el.Div{el.Class("mx-auto mt-8"),
							MetadataList(model),
						},
					},

					// slide-1
					el.Div{
						slide(1),
						Card{
							Title: el.Div{
								el.Class("text-xl font-extrabold"),
								GradientText(model.Title.Text),
							},
							Body: el.Div{
								Markdown(model.Body, el.Class("p-4 overflow-scroll max-h-[70vh]")),
							},
						}.Render(),
					},

					// slides!
					el.MapFragment(model.Comments, func(c api.Comment, idx int) el.Component {
						return el.Div{
							slide(idx + 2),
							Card{
								Badge: idx + 1,
								Title: el.A{
									el.Attr{"target", "_blank"},
									el.Href(c.Title.HRef),
									el.Code{el.Text(c.Title.Text)},
								},
								Body: el.Div{
									el.Class("overflow-scroll text-lg max-h-[70vh]"),
									Diff(c.CodeBlock),
								},
							}.Render(),
							el.Div{
								el.Class("max-w-[80%] mx-auto"),
								Markdown(c.Description),
							},
						}
					}),

					// slide-END
					el.Div{
						slide(len(model.Comments) + 2),
						el.Class("text-center font-bold"),
						el.Text("FIN"),
					},
				},
				el.Div{
					el.Class("flex-grow"),
				},
			},
			el.Script{
				el.Attrs{"type": "text/javascript"},
				el.Content{
					veun.Raw(slideJS),
				},
			},
		},
	})
}

//go:embed slideshow.js
var slideJS string
