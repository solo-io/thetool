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
  install: {{ .Install}}
{{end}}
`

	bootstrapYaml = `---
#system namespace
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: upstreams.gloo.solo.io
spec:
  group: gloo.solo.io
  names:
    kind: Upstream
    listKind: UpstreamList
    plural: upstreams
    singular: upstream
  scope: Namespaced
  version: v1
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: virtualhosts.gloo.solo.io
spec:
  group: gloo.solo.io
  names:
    kind: VirtualHost
    listKind: VirtualHostList
    plural: virtualhosts
    singular: virtualhost
  scope: Namespaced
  version: v1
---
#rbac for gloo
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: gloo-role
rules:
- apiGroups: [""]
  resources: ["pods", "services"]
  verbs: ["get", "watch", "list"]
- apiGroups: ["gloo.solo.io/v1"]
  resources: ["upstreams", "virtualhosts"]
  verbs: ["*"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: gloo-cluster-admin-binding
subjects:
- kind: ServiceAccount
  name: default
  namespace: gloo-system
roleRef:
  kind: ClusterRole
  name: gloo-role
  apiGroup: rbac.authorization.k8s.io
`
)

var (
	helmValuesTemplate *template.Template
)

func init() {
	helmValuesTemplate = template.Must(template.New("helm").Parse(helmValuesContent))
}
