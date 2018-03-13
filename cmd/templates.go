package cmd

import (
	"text/template"
)

const (
	helmValuesContent = `#values for Gloo Helm chart
ingress:
  type: NodePort
  replicaCount: 1
  port: 8080
  securePort: 8443
  adminPort: 19000
  image: "{{ .EnvoyImage }}"
  imageTag: {{ .EnvoyTag }}
  imagePullPolicy: IfNotPresent

gloo:
  type: NodePort
  replicaCount: 1
  port: 8081
  image: "{{ .GlooImage }}"
  imageTag: {{ .GlooTag }}
  imagePullPolicy: IfNotPresent

{{ $user := .DockerUser }}
#  add-ons {{ range .Addons }}
{{.SafeName}}:
  image: {{if .Image }}{{.Image}}{{else}}{{$user}}/{{.Name}}{{end}}
  imageTag: {{ .ImageTag}}
  imagePullPolicy: IfNotPresent
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
