package resources

import "embed"

//go:embed html/index.html
var Html []byte

//go:embed html/webmini.svg
var Svg []byte

//go:embed html/favicon.ico
var Favicon []byte

//go:embed html/static
var Static embed.FS
