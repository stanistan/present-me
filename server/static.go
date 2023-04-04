package presentme

import (
	"embed"
)

//go:embed static/*
var StaticContent embed.FS
