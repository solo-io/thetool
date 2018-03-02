package cmd

import (
	"text/template"
)

const (
	helmValuesContent = `#custom values for Gloo Helm chart
#gateway
gw:
  type: NodePort
  replicaCount: 1
  port: 8080
  securePort: 8443
  adminPort: 19000
  image: "{{ .EnvoyImage }}"
  imageTag: {{ .EnvoyTag }}
  serviceCluster: "envoy"
  serviceNode: "envoy"

gloo:
  type: NodePort
  replicaCount: 1
  port: 8081
  image: "{{ .GlooImage }}"
  imageTag: {{ .GlooTag }}

fdiscovery:
  type: ClusterIP
  replicaCount: 1
  port: 8080
  image: "{{ .FunctionDiscoveryImage }}"
  imageTag: "{{ .FunctionDiscoveryTag }}"
  enabled: true 
`
)

var (
	helmValuesTemplate *template.Template
)

func init() {
	helmValuesTemplate = template.Must(template.New("helm").Parse(helmValuesContent))
}
