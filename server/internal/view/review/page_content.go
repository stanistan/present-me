package review

import (
	"github.com/stanistan/present-me/internal/api"
	"github.com/stanistan/present-me/internal/github"
	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"
)

func PageContent(p github.ReviewParamsMap, model api.Review) veun.AsView {
	return Layout(p, model, LayoutParams{
		Content: el.Fragment{
			el.Class("relative"),
			el.Div{
				el.Class("h-full"),
				el.Div{
					el.Class("gap-3"),
					el.Div{
						el.Class("pt-4"),
						MetadataList(model),
					},
					Card{
						Title: el.Div{
							el.Class("text-xl font-extrabold"),
							GradientText(el.Text(model.Title.Text)),
						},
						Body: el.Div{
							el.Class("p-4"),
							Markdown(model.Body),
						},
					}.Render(),
					el.MapFragment(model.Comments, func(c api.Comment, idx int) el.Component {
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
									el.Class("p-3 flex-none md:w-2/5 text-md"),
									Markdown(c.Description),
								},
								el.Div{
									el.Class("flex-grow overflow-scroll text-sm md:border-l border-t md:border-t-0"),
									Diff(c.CodeBlock),
								},
							},
						}.Render()
					}),
				},
			},
		},
	})
}
