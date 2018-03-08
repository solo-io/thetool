package cmd

import (
	"text/template"
)

const (
	helmValuesContent = `#custom values for Gloo Helm chart
ingress:
  type: NodePort
  replicaCount: 1
  port: 8080
  securePort: 8443
  adminPort: 19000
  image: "{{ .EnvoyImage }}"
  imageTag: {{ .EnvoyTag }}

gloo:
  type: NodePort
  replicaCount: 1
  port: 8081
  image: "{{ .GlooImage }}"
  imageTag: {{ .GlooTag }}

{{ $user := .DockerUser }}
#  services {{ range .Services }}
{{.SafeName}}:
  image: {{if .Image }}{{.Image}}{{else}}{{$user}}/{{.Name}}{{end}}
  imageTag: {{ .ImageTag}}
  enable: {{ .Enable}}
  {{if .ConfigOnly}}configOnly: {{ .ConfigOnly }}{{end}}
{{end}}
`
)

var (
	helmValuesTemplate *template.Template
)

func init() {
	helmValuesTemplate = template.Must(template.New("helm").Parse(helmValuesContent))
}
