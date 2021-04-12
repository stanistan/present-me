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

---
{{ end }}

## Other Files

[1]: {{.PR.HTMLURL}}
