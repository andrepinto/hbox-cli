package configs

import (
	_ "embed"
	"text/template"
)

var (
	//go:embed templates/config.toml
	initConfigEmbed       string
	initConfigTemplate, _ = template.New("initConfig").Parse(initConfigEmbed)
)
