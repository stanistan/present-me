package review

import (
	"fmt"

	"github.com/stanistan/present-me/internal/api"
	"github.com/stanistan/present-me/internal/github"
	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"
)

func SlideContent(p github.ReviewParamsMap, model api.Review) el.Div {
	toShow := 0
	slide := func(idx int) el.AttrFunc {
		cl := "visible"
		if idx != toShow {
			cl = "hidden"
		}

		return el.Class("slide", fmt.Sprintf("slide-%d", idx), cl)
	}

	return el.Div{
		topBar(
			veun.Views{
				el.Span{
					el.Class("px-3"),
					el.Text(fmt.Sprintf("%s/%s#%s", p.Owner, p.Repo, p.Pull)),
				},
				el.Div{
					el.Class("inline-block bg-slate-50 shadow-inner text-black px-2 py-1 rounded-sm text-xs gap-3"),
					reviewLink(p, "cards"),
					el.Text(" | "),
					reviewLink(p, "slides"),
				},
			},

			el.Button{
				el.Class("text-xs px-2"),
				el.Attr{"id", "play-full-screen"},
				el.Text("▶️"),
			},
		),

		el.Div{el.Class("relative h-[95vh]"),

			el.Div{
				el.Class("bg-white flex flex-col h-full"),
				el.Attr{"id", "slideshow"},

				el.Div{el.Class("flex-grow")},
				el.Div{el.Class("flex-0 max-w-[2200px] mx-auto"),

					// slide-0
					el.Div{
						slide(0),
						el.Div{el.Class("text-6xl font-extrabold text-center"),
							GradientText(el.Text(model.Title.Text)),
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
								GradientText(el.Text(model.Title.Text)),
							},
							Body: el.Div{
								el.Class("p-4"),
								Markdown(model.Body),
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
								el.Class("max-w-[80%] mx-auto markdown"),
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
		},
		el.Script{
			el.Attrs{"type": "text/javascript"},
			el.Content{
				veun.Raw(slideJS),
			},
		},
	}
}

const slideJS = `
document.getElementById("play-full-screen").addEventListener("click", () => {
	document.getElementById("slideshow").requestFullscreen();
});

window.addEventListener("keyup", (e) => {
	if (e.defaultPrevented) {
		return;
	}

	var action = "";
	switch (e.key) {
	case "ArrowLeft":
		action = "prev";
		break;
	case "ArrowRight":
		action = "next";
		break;
	default:
		return;
	}

	e.preventDefault();

	var currentSlide = document.querySelectorAll(".slide.visible")[0];
	var nextSlide = null;
	if (action == "prev") {
		nextSlide = currentSlide.previousElementSibling;
	} else {
		nextSlide = currentSlide.nextElementSibling;
	}

	if (nextSlide != null) {
		var currentCl = currentSlide.classList;
		currentCl.remove("visible");
		currentCl.add("hidden");

		var nextCl = nextSlide.classList;
		nextCl.add("visible");
		nextCl.remove("hidden");
	}
})
`
