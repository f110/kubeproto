.PHONY: build
build: pkg/apis/github.generated.deepcopy.go \
		pkg/apis/github.generated.register.go

pkg/apis/github.generated.deepcopy.go: pkg/apis/github.proto
	bazel build //example/pkg/apis:github_proto_deepcopy
	@rm -f $@
	cp ../bazel-bin/example/pkg/apis/github_proto_deepcopy_/go.f110.dev/kubeproto/example/pkg/apis/github.generated.deepcopy.go $(@D)
	@chmod 644 $@

pkg/apis/github.generated.register.go: pkg/apis/github.proto
	bazel build //example/pkg/apis:github_proto_register
	@rm -f $@
	cp ../bazel-bin/example/pkg/apis/github_proto_register_/go.f110.dev/kubeproto/example/pkg/apis/github.generated.register.go $(@D)
	@chmod 644 $@

.PHONY: clean
clean:
	rm -f pkg/apis/github.generated.deepcopy.go
	rm -f pkg/apis/github.generated.register.go