jobs: test_all: {
	command: "test"
	targets: ["//..."]
	all_revision: true
	github_status: true
	platforms: ["@rules_go//go/toolchain:linux_amd64"]
	cpu_limit: "2000m"
	event: ["push", "pull_request"]
}
