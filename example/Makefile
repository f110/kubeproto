.PHONY: build
build: githubv1alpha1 githubv1alpha2 \
		crd/github.crd.yaml

.PHONY: githubv1alpha1
githubv1alpha1: pkg/apis/githubv1alpha1/github.proto
	bazel build //example/pkg/apis/githubv1alpha1:github_proto_k8s
	@rm -f pkg/apis/$@/github.generated.deepcopy.go pkg/apis/$@/github.generated.register.go
	cp ../bazel-bin/example/pkg/apis/$@/github_proto_k8s_/go.f110.dev/kubeproto/example/pkg/apis/$@/github.generated.deepcopy.go pkg/apis/$@
	cp ../bazel-bin/example/pkg/apis/$@/github_proto_k8s_/go.f110.dev/kubeproto/example/pkg/apis/$@/github.generated.register.go pkg/apis/$@

.PHONY: githubv1alpha2
githubv1alpha2: pkg/apis/githubv1alpha2/github.proto
	bazel build //example/pkg/apis/githubv1alpha2:github_proto_k8s
	@rm -f pkg/apis/$@/github.generated.deepcopy.go pkg/apis/$@/github.generated.register.go
	cp ../bazel-bin/example/pkg/apis/$@/github_proto_k8s_/go.f110.dev/kubeproto/example/pkg/apis/$@/github.generated.deepcopy.go pkg/apis/$@
	cp ../bazel-bin/example/pkg/apis/$@/github_proto_k8s_/go.f110.dev/kubeproto/example/pkg/apis/$@/github.generated.register.go pkg/apis/$@

pkg/client/github.generated.client.go: pkg/apis/githubv1alpha1/github.proto pkg/apis/githubv1alpha2/github.proto
	bazel build //example/pkg/client:github
	@rm -f $@
	cp ../bazel-bin/example/pkg/client/github.generated.client.go $(@D)
	@chmod 644 $@

crd/github.crd.yaml: pkg/apis/githubv1alpha1/github.proto pkg/apis/githubv1alpha2/github.proto
	bazel build //example/crd:github
	@rm -f $@
	cp ../bazel-bin/example/crd/github.crd.yaml $(@D)
	@chmod 644 $@

.PHONY: clean
clean:
	rm -f pkg/apis/githubv1alpha1/github.generated.deepcopy.go
	rm -f pkg/apis/githubv1alpha1/github.generated.register.go
	rm -f pkg/apis/githubv1alpha2/github.generated.deepcopy.go
	rm -f pkg/apis/githubv1alpha2/github.generated.register.go
	rm -f crd/github.crd.yaml