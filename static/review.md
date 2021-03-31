# ([#{{.Params.Number}}][1]) {{.PR.Title}}

| | |
|:-|-:|
| Source | {{ .PR.HTMLURL }} |
| Author | [{{.Review.User.Login}}]({{.Review.User.HTMLURL}}) |

---

{{ .PR.Body | safe }}

{{ .Review.Body | safe }}

---

{{ range .Comments }}
## `{{ .Path }}`

```diff
{{.DiffHunk | safe}}
```

{{.Body | stripLeadingNumber }}

---
{{ end }}

[1]: {{.PR.HTMLURL}}
