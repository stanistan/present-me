package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"

	"github.com/stanistan/crap"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{owner}/{repo}/pull/{number}/{reviewID}", renderReview).
		Methods("GET").
		Name("review")

	s := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}

func renderReview(w http.ResponseWriter, r *http.Request) {
	params, err := crap.ReviewParamsFromMap(mux.Vars(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := crap.Run(crap.Context{
		Ctx:    r.Context(),
		Client: github.NewClient(nil),
	}, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithCodeBlockOptions(func(_ highlighting.CodeBlockContext) []html.Option {
					return []html.Option{}
				}),
			),
		),
	)

	var buf bytes.Buffer
	err = md.Convert(b.Bytes(), &buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = bodyContentTemplate.Execute(w, map[string]interface{}{
		"Body": template.HTML(buf.Bytes()),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

var bodyContentTemplate = template.Must(template.New("content").Parse(`
<html>
	<head>
		<style type="text/css">
			body{border-top:5px solid #f6cde0;background:#fcfcfc;color:#333;font-family:Cochin, Times, serif;font-size:1.3em;line-height:1.55em;margin:0;padding:0}.content{margin:8em auto 0;max-width:33em;padding:2em;position:relative}article{position:relative}h1{font-size:1.5em;font-weight:bold;margin:0;padding:1.5em 0 0}h2{font-size:1.25em;font-weight:bold;margin:0;padding:1.25em 0 0}h3{font-size:1.15em;font-weight:bold;margin:0;padding:1.15em 0 0}h4{font-size:1em;font-weight:bold;margin:0;padding:1em 0 0}hr{border:0 dotted #f6cde0;border-top-width:1px;height:0}pre,.mono,code{font-family:monospace}.small{font-size:75%}.smaller{font-size:85%}.faded{color:#999}.center{text-align:center}a,.link-like{color:#0000a5;font-weight:bolder}a:hover,.link-like:hover{color:#2525ff;text-decoration:none}ol,ul{list-style-position:outside;margin:0;margin-left:1.6em;padding:0}ol li,ul li{margin:0 0 0.1em 0;padding:0 0 0 0.2em}.outside{margin-left:-2em}ul.outside,ol.outside{margin-left:-0.1em}code{background:#fdf4f8;border-bottom:1px dotted #f6cde0;padding:0.3em;font-size:0.9em}pre,blockquote{background:white;line-height:1.45em;margin-left:-2em;margin-right:-2em}pre{font-size:0.85em;border-bottom:1px dotted #f6cde0;padding:1em 2em;overflow-x:scroll}blockquote{font-size:0.9em;font-style:italic;padding:0.5em 2em;border-left:1px dotted #f6cde0}.tag{margin-right:0.3em;border-radius:10px;color:#333;text-decoration:none}.tag:hover{color:#000;text-decoration:underline}.toc{width:15em;font-size:0.8em;padding:1em;background:#fff;line-height:1.45;border-bottom:1px dotted #f6cde0}.toc a{color:#333;font-weight:normal}.anchor-link{position:absolute;margin-left:-1em;text-decoration:none}article .img-container{margin:0.2em -2em 0}article .img-container img{width:100%;height:auto}@media only screen and (max-width: 500px){body{font-size:1.15em}.content{padding:1.8em}.smaller{font-size:90%}.small{font-size:83%}.outside{margin-left:inherit}article .img-container{margin:0.2em -1em 0}}dl{margin:1.5em 0}dl dt{clear:left;float:left;text-decoration:underline;text-align:right;width:5.5em}dl dt::after{content:":"}dl dd{display:block;float:left}.clear{clear:both}.invoice-heading{padding:1.5em 0}.invoice-heading h1{float:left;padding:0}.invoice-heading .invoice-heading-code{float:right;font-weight:bold}table{padding-top:1.5em;width:100%}table th,table td{text-align:left}table th.right,table td.right{text-align:right}table tfoot td{font-weight:bold}table tfoot tr:first-child td{border-top:1px dotted #f6cde0;padding-top:0.3em}table tbody tr:last-child td{padding-bottom:0.3em}table tbody tr:first-child td{padding-top:0.3em}table thead tr:last-child th{padding-bottom:0.3em;border-bottom:1px dotted #ccc}details{font-size:85%}details summary{font-weight:bold;color:#0000a5}details summary:hover{color:#2525ff;cursor:pointer}
		</style>
	</head>
	<body>
		<div class="content">
			<article>
			{{ .Body }}
			</article>
		</div>
	</body>
</html>
`))
