package github

import "path"

func detectLanguage(p string) string {
	var (
		base = path.Base(p)
		ext  = path.Ext(p)
	)

	switch ext {
	case "":
		if base == "Dockerfile" {
			return "docker"
		}
		return "bash"
	case ".rs":
		return "rust"
	case ".vue":
		return "html"
	case ".Dockerfile":
		return "docker"
	}

	return ext[1:]
}
