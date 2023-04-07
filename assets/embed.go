package assets

import (
	_ "embed"
	"text/template"
)

//go:embed templates/cloud-init.tmpl
var cloudInitTemplateRaw []byte
var CloudInitTemplate *template.Template

func init() {
	CloudInitTemplate = template.Must(template.New("cloudinit").Parse(string(cloudInitTemplateRaw)))
}
