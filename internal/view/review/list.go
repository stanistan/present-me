package review

import (
	"fmt"
	"sort"

	"github.com/stanistan/present-me/internal/api"
	"github.com/stanistan/present-me/internal/github"
	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"
)

func SourcesList(p github.ReviewParamsMap, m api.Review) veun.AsView {
	var tagCounts = tagsMeta(m)
	var ks []string
	for tag, _ := range tagCounts {
		ks = append(ks, tag)
	}

	sort.Strings(ks)

	var content []el.Param
	for _, k := range ks {
		var tag string
		if k != "" {
			tag = "#" + k
		} else {
			tag = "/unlabelled/"
		}

		content = append(content, el.Li{
			el.A{
				el.Class(
					"block",
					"px-3", "py-2", "border-b",
					"text-sm", "bg-gray-white hover:bg-pink-50",
					"font-mono", "hover:text-indigo-900",
				),
				el.Href(fmt.Sprintf("/%s/%s/pull/%s/tag%s/cards", p.Owner, p.Repo, p.Pull, k)),
				el.Div{
					el.Class("underline"),
					el.Text(tag),
				},
				el.Div{
					el.Class("text-xs"),
					el.Text(fmt.Sprintf("%d", tagCounts[k])),
					el.Text(" comment(s)"),
				},
			},
		})
	}

	return Layout(p, m, LayoutParams{
		Content: el.Fragment{
			el.Class("max-w-96", "p-3", "mx-auto"),
			Card{
				Title: el.Div{
					el.Class("text-md", "text-center"),
					el.Text("Choose one:"),
				},
				Body: el.Div{
					el.Ol{
						el.Class("list-decimal"),
						el.Fragment(content),
					},
				},
			}.Render(),
		},
	})
}

func tagsMeta(m api.Review) map[string]int {
	out, ok := m.MetaData["reviewTags"]
	if !ok {
		return nil
	}

	tags, ok := out.(map[string]int)
	if !ok {
		return nil
	}

	return tags
}
