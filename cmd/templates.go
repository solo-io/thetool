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
  port: 80
  image: "{{ .EnvoyImage }}"
  imageTag: {{ .EnvoyTag }}
  serviceCluster: cluster
  serviceNode: node

gloo:
  type: NodePort
  replicaCount: 1
  port: 80
  image: "{{ .GlooImage }}"
  imageTag: {{ .GlooTag }}

fdiscovery:
  type: ClusterIP
  replicaCount: 1
  port: 8080
  image: "soloio/glue-discovery"
  imageTag: "9cf492d"
  enabled: true 
`
)

var (
	helmValuesTemplate *template.Template
)

func init() {
	helmValuesTemplate = template.Must(template.New("helm").Parse(helmValuesContent))
}
