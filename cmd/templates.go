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

function_discovery:
  type: ClusterIP
  replicaCount: 1
  port: 8080
  image: "{{ .FunctionDiscoveryImage }}"
  imageTag: "{{ .FunctionDiscoveryTag }}"
  enabled: true 

# features
{{ range .Features }}
{{.Name}}_enabled: {{.Enabled}}
{{end}}
`
)

var (
	helmValuesTemplate *template.Template
)

func init() {
	helmValuesTemplate = template.Must(template.New("helm").Parse(helmValuesContent))
}
