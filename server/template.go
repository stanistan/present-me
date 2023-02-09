package presentme

import (
	"bytes"
	_ "embed"
	"html/template"
	"io"

	"github.com/rs/zerolog/log"
)

func ErrorPage(w io.Writer, err *Error) error {
	return errorTemplate.Execute(w, struct {
		Error *Error
	}{
		Error: err,
	})
}

func IndexPage(w io.Writer, url, err string) error {
	return indexTemplate.Execute(w, struct {
		URL, Err string
		ReadMe   string
	}{
		URL:    url,
		Err:    err,
		ReadMe: readme,
	})
}

func intoTemplate(w io.Writer, title string, bytes []byte) error {
	return htmlTemplate.Execute(w, struct {
		Body  template.HTML
		Title string
	}{
		Body:  template.HTML(bytes),
		Title: title,
	})
}

func reviewBody(w io.Writer, r *ReviewModel) error {
	return reviewTemplate.Execute(w, r)
}

func asSlide(w io.Writer, title string, bytes []byte) error {
	return slideTemplate.Execute(w, struct {
		Body  template.HTML
		Title string
	}{
		Body:  template.HTML(bytes),
		Title: title,
	})
}

var templateFuncMap = template.FuncMap{
	"render_md": func(s string) template.HTML {
		var buf bytes.Buffer
		err := intoHTML(&buf, []byte(s))
		if err != nil {
			log.Fatal().Err(err).Msg("")
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
	//go:embed README.md
	readme string

	//go:embed templates/html.html
	htmlBytes string

	//go:embed templates/review.md
	reviewBytes string

	//go:embed templates/slides.html
	slideBytes string

	//go:embed templates/index.html
	indexBytes string

	//go:embed templates/error.html
	errorBytes string

	reviewTemplate = templateMust("review", reviewBytes)
	htmlTemplate   = templateMust("html", htmlBytes)
	slideTemplate  = templateMust("slide", slideBytes)
	indexTemplate  = templateMust("index", indexBytes)
	errorTemplate  = templateMust("error", errorBytes)
)
