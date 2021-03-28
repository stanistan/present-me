# (#{{.Params.Number}}) {{.PR.Title}}

{{ .Review.Body }}

---

{{ range .Comments }}
### `{{ .Path }}`

```diff
{{.DiffHunk}}
```

{{.Body | stripLeadingNumber }}

---
{{ end }}
