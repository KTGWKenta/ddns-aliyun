// +build go1.16

package config

import _ "embed"

//go:embed template.cue
var innerTemplate string
