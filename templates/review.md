{{ $r := . }}
# ([#{{.Params.Number}}][1]) {{.PR.Title}}

{{ define "BASE" }}/{{.Params.Owner}}/{{.Params.Repo}}/pull/{{.Params.Number}}/{{.Params.ReviewID}}{{end}}

| | |
|:-|-:|
| *Author* | [{{.Review.User.Login}}]({{.Review.User.HTMLURL}}) |
| *Source* | {{ .PR.HTMLURL }} |
| *Post*   | [{{template "BASE" . }}/post](post) |
| *Slides* | [{{template "BASE" . }}/slides](slides) |
| *.md*    | [{{template "BASE" . }}/md](md) |


{{ .PR.Body | safe }}

---

## Notable Changes

{{ .Review.Body | safe }}

---

{{ range .Comments }}
### `{{ .Path }}` ([thread]({{.HTMLURL}}))

```diff
{{.DiffHunk | safe}}
```

{{.Body | stripLeadingNumber }}

{{ with $f := $r.CommitFile .Path }}
Changes:
: **+{{$f.Additions}}/-{{$f.Deletions}}** {{$f.Status}}
{{ end }}

---
{{ end }}

## Other Files

| File | +/- | |
|:-|-:|:-|{{ range $idx, $f := .Files }}{{ if not $f.IsAnnotated }}
| [`{{ $f.File.Filename }}`]({{$f.File.RawURL}}) | **+{{$f.File.Additions}}/-{{$f.File.Deletions}}** | {{$f.File.Status}} | {{end}}{{ end }}

[1]: {{.PR.HTMLURL}}
