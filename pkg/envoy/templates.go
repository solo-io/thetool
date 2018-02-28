package envoy

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/solo-io/thetool/pkg/downloader"
	"github.com/solo-io/thetool/pkg/feature"
)

const (
	buildContent = `package(default_visibility = ["//visibility:public"])

load(
    "@envoy//bazel:envoy_build_system.bzl",
    "envoy_cc_binary",
    "envoy_cc_library",
    "envoy_cc_test",
)

envoy_cc_binary(
    name = "envoy",
    repository = "@envoy",
    deps = [{{range .}}
		"@{{.Name}}//:filter_lib",{{end}}
		"@envoy//source/exe:envoy_main_entry_lib",
    ],
)
`

	workspaceContent = `load('@bazel_tools//tools/build_defs/repo:git.bzl', 'git_repository')
{{range .}}
local_repository(
    name = "{{.Name}}",
    path = "{{. | path}}",
)

{{end}}
bind(
    name = "boringssl_crypto",
    actual = "//external:ssl",
)

ENVOY_COMMON_SHA = "0fefba1a893e9d8e24700b0c5919f301a3d10249"  # Feb 19, 2018 (let the child filter  retrieve the function.)

# load solo common
git_repository(
   name = "envoy_common",
   remote = "git@github.com:solo-io/envoy-common",
   commit = ENVOY_COMMON_SHA,
)

# some dependencies that are hard coded for now; need to fix
JSON_SHA = "c8ea63a31bbcf652d61490b0ccd86771538f8c6b"

new_http_archive(
   name = "json",
   strip_prefix = "json-" + JSON_SHA + "/single_include/nlohmann",
   url = "https://github.com/nlohmann/json/archive/" + JSON_SHA + ".zip",
   build_file_content = """
cc_library(
   name = "json-lib",
   hdrs = ["json.hpp"],
   visibility = ["//visibility:public"],
)
   """
)


INJA_SHA = "74ad4281edd4ceca658888602af74bf2050107f0"

new_http_archive(
   name = "inja",
   strip_prefix = "inja-" + INJA_SHA + "/src",
   url = "https://github.com/pantor/inja/archive/" + INJA_SHA + ".zip",
   build_file_content = """
cc_library(
   name = "inja-lib",
   hdrs = ["inja.hpp"],
   visibility = ["//visibility:public"],
)
   """
)

http_archive(
    name = "envoy",
    strip_prefix = "envoy-{{ envoyHash }}",
    url = "https://github.com/envoyproxy/envoy/archive/{{ envoyHash }}.zip",
)

load("@envoy//bazel:repositories.bzl", "envoy_dependencies")
envoy_dependencies()
load("@envoy//bazel:cc_configure.bzl", "cc_configure")
cc_configure()
load("@envoy_api//bazel:repositories.bzl", "api_dependencies")
api_dependencies()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
load("@com_lyft_protoc_gen_validate//bazel:go_proto_library.bzl", "go_proto_repositories")
go_proto_repositories(shared=0)
go_rules_dependencies()
go_register_toolchains()
load("@io_bazel_rules_go//proto:def.bzl", "proto_register_toolchains")
proto_register_toolchains()
`

	dockerfile = `FROM ubuntu:16.04

ADD WORKSPACE /etc/envoy.WORKSPACE
ADD envoy /usr/local/bin/envoy

CMD /usr/local/bin/envoy -c /etc/envoy.yaml --service-cluster $CLUSTER --service-node $NODE`
)

var (
	buildTemplate     *template.Template
	workspaceTemplate *template.Template

	envoyHash = "29989a38c017d3be5aa3c735a797fcf58b754fe5"
	workDir   = "repositories"
)

func init() {
	buildTemplate = template.Must(template.New("build").Parse(buildContent))
	funcMap := template.FuncMap{
		"path":      path,
		"envoyHash": func() string { return envoyHash },
	}
	workspaceTemplate = template.Must(template.New("workspace").
		Funcs(funcMap).Parse(workspaceContent))
}

func path(f feature.Feature) string {
	if strings.HasSuffix(f.Repository, ".git") {
		return fmt.Sprintf("%s/%s/%s", workDir, downloader.RepoDir(f.Repository), f.EnvoyDir)
	}

	if isGitHubHTTP(f.Repository) {
		return fmt.Sprintf("%s/%s-%s/envoy", workDir, f.Name, f.Revision)
	}

	return fmt.Sprintf("%s/%s/%s", workDir, downloader.RepoDir(f.Repository), f.EnvoyDir)
}

func isGitHubHTTP(url string) bool {
	return strings.HasPrefix(url, "https://github.com/")
}
