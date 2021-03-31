package presentme

import (
	"bytes"
	_ "embed"
	"html/template"
	"io"
	"log"
)

func intoTemplate(w io.Writer, bytes []byte) error {
	return htmlTemplate.Execute(w, struct {
		Body template.HTML
	}{
		Body: template.HTML(bytes),
	})
}

func reviewBody(w io.Writer, r *ReviewModel) error {
	return reviewTemplate.Execute(w, r)
}

var templateFuncMap = template.FuncMap{
	"render_md": func(s string) template.HTML {
		var buf bytes.Buffer
		err := intoHTML(&buf, []byte(s))
		if err != nil {
			log.Fatal(err)
		}
		return template.HTML(buf.Bytes())
	},
	"stripLeadingNumber": func(s string) template.HTML {
		return template.HTML(stripLeadingNumber(s))
	},
	"safe": func(s string) template.HTML {
		return template.HTML(s)
	},
}

func templateMust(n, content string) *template.Template {
	return template.Must(template.New(n).
		Funcs(templateFuncMap).
		Parse(content),
	)
}

var (
	//go:embed static/html.html
	htmlBytes string

	//go:embed static/review.md
	reviewBytes string

	reviewTemplate = templateMust("review", reviewBytes)
	htmlTemplate   = templateMust("html", htmlBytes)
)
