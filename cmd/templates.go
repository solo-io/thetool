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

control_plane:
  replicaCount: 1
  port: 8081
  image: "{{ .GlooImage }}"
  imageTag: "{{ .GlooTag }}"
  imagePullPolicy: IfNotPresent

{{ $user := .DockerUser }} {{ $glooTag := .GlooTag }}
#  add-ons {{ range .Addons }}
{{.SafeName}}:
  {{if .IsGlooAddon }}image: "{{$user}}/{{.Name}}"
  imageTag: "{{$glooTag}}"
  {{end}}imagePullPolicy: IfNotPresent{{range $k, $v := .Configuration }}
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
