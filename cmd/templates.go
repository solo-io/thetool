package cmd

import (
	"text/template"
)

const (
	helmValuesContent = `#custom values for Glue Helm chart
#gateway
gw:
  type: NodePort
  replicaCount: 1
  port: 80
  image: "{{ .EnvoyImage }}"
  imageTag: {{ .EnvoyTag }}
  serviceCluster: cluster
  serviceNode: node

glue:
  type: NodePort
  replicaCount: 1
  port: 80
  image: "{{ .GlueImage }}"
  imageTag: {{ .GlueTag }}"
`
)

var (
	helmValuesTemplate *template.Template
)

func init() {
	helmValuesTemplate = template.Must(template.New("helm").Parse(helmValuesContent))
}
