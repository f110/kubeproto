.PHONY: build
build: pkg/apis/githubv1alpha1/github.generated.deepcopy.go \
		pkg/apis/githubv1alpha1/github.generated.register.go \
		pkg/apis/githubv1alpha2/github.generated.deepcopy.go \
		pkg/apis/githubv1alpha2/github.generated.register.go \
		crd/github.crd.yaml

pkg/apis/githubv1alpha1/github.generated.deepcopy.go: pkg/apis/githubv1alpha1/github.proto
	bazel build //example/pkg/apis/githubv1alpha1:github_proto_deepcopy
	@rm -f $@
	cp ../bazel-bin/example/pkg/apis/githubv1alpha1/github_proto_deepcopy_/go.f110.dev/kubeproto/example/pkg/apis/githubv1alpha1/github.generated.deepcopy.go $(@D)
	@chmod 644 $@

pkg/apis/githubv1alpha1/github.generated.register.go: pkg/apis/githubv1alpha1/github.proto
	bazel build //example/pkg/apis/githubv1alpha1:github_proto_register
	@rm -f $@
	cp ../bazel-bin/example/pkg/apis/githubv1alpha1/github_proto_register_/go.f110.dev/kubeproto/example/pkg/apis/githubv1alpha1/github.generated.register.go $(@D)
	@chmod 644 $@

pkg/apis/githubv1alpha2/github.generated.deepcopy.go: pkg/apis/githubv1alpha2/github.proto
	bazel build //example/pkg/apis/githubv1alpha2:github_proto_deepcopy
	@rm -f $@
	cp ../bazel-bin/example/pkg/apis/githubv1alpha2/github_proto_deepcopy_/go.f110.dev/kubeproto/example/pkg/apis/githubv1alpha2/github.generated.deepcopy.go $(@D)
	@chmod 644 $@

pkg/apis/githubv1alpha2/github.generated.register.go: pkg/apis/githubv1alpha2/github.proto
	bazel build //example/pkg/apis/githubv1alpha2:github_proto_register
	@rm -f $@
	cp ../bazel-bin/example/pkg/apis/githubv1alpha2/github_proto_register_/go.f110.dev/kubeproto/example/pkg/apis/githubv1alpha2/github.generated.register.go $(@D)
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