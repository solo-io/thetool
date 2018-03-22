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
  imageTag: "{{ .EnvoyTag }}"
  imagePullPolicy: IfNotPresent

gloo:
  replicaCount: 1
  port: 8081
  image: "{{ .GlooImage }}"
  imageTag: "{{ .GlooTag }}"
  imagePullPolicy: IfNotPresent

{{ $user := .DockerUser }}
#  add-ons {{ range .Addons }}
{{.SafeName}}:
  {{if .Repository }}image: "{{$user}}/{{.Name}}"
  imageTag: "{{ .ImageTag }}"
  {{end}}imagePullPolicy: IfNotPresent
  enable: {{ .Enable}}{{range $k, $v := .Configuration }}
  {{$k}}: {{$v}}{{end}}
{{end}}
`
)

var (
	helmValuesTemplate *template.Template
)

func init() {
	helmValuesTemplate = template.Must(template.New("helm").Parse(helmValuesContent))
}
