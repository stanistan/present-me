package review

import (
	"fmt"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"

	"github.com/stanistan/present-me/internal/api"
	"github.com/stanistan/present-me/internal/github"
)

func flexNone(fs ...el.Param) el.Div {
	return el.Div{el.Class("flex-none"), el.Fragment(fs)}
}

func flexSpace() el.Div {
	return el.Div{el.Class("flex-grow")}
}

func flexRow(fs ...el.Param) el.Div {
	return el.Div{el.Class("flex", "flex-row"), el.Fragment(fs)}
}

func topBar(main, right veun.AsView) el.Div {
	return el.Div{
		el.Class(
			"py-2",
			"bg-gradient-to-b from-gray-800 to-black",
			"text-white font-mono text-sm",
		),
		flexRow(
			flexNone(el.A{
				el.Href("/"),
				el.Class("font-bold hover:text-pink-300 ml-3"),
				el.Text(" / "),
			}),
			flexSpace(),
			flexNone(el.Content{main}),
			flexSpace(),
			flexNone(el.Class("pr-3"), el.Content{right}),
		),
	}
}

func GradientText(text string) el.Span {
	return el.Span{
		el.Class("bg-clip-text text-transparent bg-gradient-to-r from-pink-600 to-violet-900"),
		el.Text(text),
	}
}

func reviewLink(p github.ReviewParamsMap, to string) el.A {
	href := fmt.Sprintf("/%s/%s/pull/%s/%s/%s", p.Owner, p.Repo, p.Pull, p.Source(), to)
	class := "underline"
	if p.Kind == to {
		class = "font-bold no-underline"
	}

	return el.A{
		el.Href(href),
		el.Class(class, "hover:no-underline"),
		el.Text(to),
	}
}

func Markdown(in string) el.Component {
	if in == "" {
		return nil
	}

	return el.Div{
		el.Class("markdown"),
		el.Content{markdown(in)},
	}
}

func Diff(code api.CodeBlock) el.Pre {
	lang := "language-" + code.Language
	if code.IsDiff {
		lang = "language-diff-" + code.Language
	}

	return el.Pre{
		el.Class("bg-gray-50", lang),
		el.Code{
			el.Class("diff-highlight", lang),
			el.Text(code.Content),
		},
	}
}

type Card struct {
	Badge       int
	Title, Body veun.AsView
}

func (c Card) Render() el.Component {
	return el.Div{
		el.Class("m-4", "border border-slate-300", "rounded-xl", "overflow-hidden", "shadow"),
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

func MetadataList(m api.Review) el.Div {
	return el.Div{
		el.MapFragment(m.Links, func(l api.LabelledLink, _ int) el.Div {
			return el.Div{
				el.Class("grid grid-cols-2 gap-4 text-xs font-mono"),
				el.Div{
					el.Class("text-right p-1"),
					el.A{
						el.Class("underline hover:no-underline"),
						el.Href(l.HRef),
						el.Text(l.Text),
					},
				},
				el.Div{
					el.Class("p-1"),
					el.Strong{
						el.Text(l.Label),
					},
				},
			}
		}),
	}
}
