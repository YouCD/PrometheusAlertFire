package dist

import "embed"

//go:embed  * css/* js/*
var Dist embed.FS
